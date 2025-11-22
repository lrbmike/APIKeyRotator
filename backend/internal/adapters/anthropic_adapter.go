package adapters

import (
	"fmt"
	"io"
	"strings"

	"api-key-rotator/backend/internal/cache"
	"api-key-rotator/backend/internal/config"
	"api-key-rotator/backend/internal/logger"
	"api-key-rotator/backend/internal/models"
	"api-key-rotator/backend/internal/services"
	"api-key-rotator/backend/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AnthropicAdapter is an adapter for the Anthropic API.
type AnthropicAdapter struct {
	*BaseLLMAdapter
}

// NewAnthropicAdapter creates a new AnthropicAdapter instance.
func NewAnthropicAdapter(cfg *config.Config, db *gorm.DB, cacheClient cache.CacheInterface,
	c *gin.Context, proxyConfig *models.ProxyConfig, action string) *AnthropicAdapter {
	return &AnthropicAdapter{
		BaseLLMAdapter: NewBaseLLMAdapter(cfg, db, cacheClient, c, proxyConfig, action),
	}
}

// ProcessRequest handles the request for the Anthropic API.
func (a *AnthropicAdapter) ProcessRequest() (*services.TargetRequest, error) {
	// 1. Authenticate the proxy request (hijack the 'x-api-key' header).
	proxyKey := a.c.GetHeader("x-api-key")

	validKeys := a.cfg.GetGlobalProxyKeys()
	isValidKey := false
	for _, key := range validKeys {
		if proxyKey == key {
			isValidKey = true
			break
		}
	}
	if !isValidKey {
		return nil, fmt.Errorf("invalid Proxy Key. Provide it via the 'x-api-key' header")
	}

	// 2. Rotate the upstream key.
	upstreamKey, err := a.RotateUpstreamKey()
	if err != nil {
		return nil, err
	}

	// 3. Build the target request.
	headers := utils.FilterRequestHeaders(a.c.Request.Header, []string{"x-api-key"})

	// 优先使用数据库中为该proxyConfig保存的APIKeyName, 否则回退到默认值
	keyName := "x-api-key"
	if a.proxyConfig.APIKeyName != nil && *a.proxyConfig.APIKeyName != "" {
		keyName = *a.proxyConfig.APIKeyName
	}

	// 优先使用数据库中为该proxyConfig保存的APIKeyLocation, 否则回退到默认值
	keyLocation := "header"
	if a.proxyConfig.APIKeyLocation != nil && *a.proxyConfig.APIKeyLocation != "" {
		keyLocation = *a.proxyConfig.APIKeyLocation
	}

	// 根据keyLocation将key添加到headers或params
	if keyLocation == "header" {
		// Inject the real upstream key and the fixed anthropic-version.
		headers[keyName] = upstreamKey
	} else if keyLocation == "query" {
		// 将key添加到查询参数
		if a.c.Request.URL.RawQuery == "" {
			a.c.Request.URL.RawQuery = fmt.Sprintf("%s=%s", keyName, upstreamKey)
		} else {
			a.c.Request.URL.RawQuery += fmt.Sprintf("&%s=%s", keyName, upstreamKey)
		}
	}
	headers["anthropic-version"] = "2023-06-01"

	baseURL := ""
	if a.proxyConfig.TargetBaseURL != nil {
		baseURL = strings.TrimSuffix(*a.proxyConfig.TargetBaseURL, "/")
	}
	finalURL := fmt.Sprintf("%s/%s", baseURL, a.action)

	params := make(map[string]string)
	for key, values := range a.c.Request.URL.Query() {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}

	body, err := io.ReadAll(a.c.Request.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}

	logger.Infof("%s: Rotated to upstream key (masked): %s", a.logPrefix, utils.MaskAPIKeyDefault(upstreamKey))

	return &services.TargetRequest{
		Method:  a.c.Request.Method,
		URL:     finalURL,
		Headers: headers,
		Params:  params,
		Body:    body,
	}, nil
}