package adapters

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"api-key-rotator/backend/internal/config"
	"api-key-rotator/backend/internal/logger"
	"api-key-rotator/backend/internal/models"
	"api-key-rotator/backend/internal/services"
	"api-key-rotator/backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// OpenAIAdapter 适配器，用于处理所有兼容OpenAI格式的API
type OpenAIAdapter struct {
	*BaseLLMAdapter
}

// NewOpenAIAdapter 创建OpenAI适配器实例
func NewOpenAIAdapter(cfg *config.Config, db *gorm.DB, redisClient *redis.Client,
	c *gin.Context, proxyConfig *models.ProxyConfig, action string) *OpenAIAdapter {
	return &OpenAIAdapter{
		BaseLLMAdapter: NewBaseLLMAdapter(cfg, db, redisClient, c, proxyConfig, action),
	}
}

// ProcessRequest 处理OpenAI格式的请求
func (a *OpenAIAdapter) ProcessRequest() (*services.TargetRequest, error) {
	// 1. 代理访问认证 (劫持官方流程)
	bodyBytes, err := io.ReadAll(a.c.Request.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}

	authSuccessful := false

	// 模式A: Header认证 (优先)
	authHeader := a.c.GetHeader("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		token := strings.TrimPrefix(authHeader, "Bearer ")
		validKeys := a.cfg.GetGlobalProxyKeys()
		for _, key := range validKeys {
			if token == key {
				authSuccessful = true
				break
			}
		}
	}

	// 模式B: Body认证 (备选)
	if !authSuccessful && len(bodyBytes) > 0 {
		var bodyJSON map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &bodyJSON); err == nil {
			if apiKey, exists := bodyJSON["api_key"]; exists {
				if apiKeyStr, ok := apiKey.(string); ok {
					validKeys := a.cfg.GetGlobalProxyKeys()
					for _, key := range validKeys {
						if apiKeyStr == key {
							authSuccessful = true
							break
						}
					}
					// 从body中移除api_key字段
					delete(bodyJSON, "api_key")
					if newBodyBytes, err := json.Marshal(bodyJSON); err == nil {
						bodyBytes = newBodyBytes
					}
				}
			}
		}
	}

	if !authSuccessful {
		return nil, fmt.Errorf("invalid Proxy Key. Provide it via 'Authorization: Bearer <key>' header or 'api_key' in JSON body")
	}

	// 2. 轮询上游密钥
	upstreamKey, err := a.RotateUpstreamKey()
	if err != nil {
		return nil, err
	}

	// 3. 构建目标请求 (偷梁换柱)
	headers := utils.FilterRequestHeaders(a.c.Request.Header, []string{"authorization"})

	// 优先使用数据库中为该proxyConfig保存的APIKeyName, 否则回退到默认值
	keyName := "Authorization"
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
		// 添加真实的Authorization头
		headers[keyName] = fmt.Sprintf("Bearer %s", upstreamKey)
	} else if keyLocation == "query" {
		// 将key添加到查询参数
		if a.c.Request.URL.RawQuery == "" {
			a.c.Request.URL.RawQuery = fmt.Sprintf("%s=%s", keyName, upstreamKey)
		} else {
			a.c.Request.URL.RawQuery += fmt.Sprintf("&%s=%s", keyName, upstreamKey)
		}
	}

	// 从数据库获取配置好的基础URL，并移除末尾可能存在的斜杠
	baseURL := ""
	if a.proxyConfig.TargetBaseURL != nil {
		baseURL = strings.TrimSuffix(*a.proxyConfig.TargetBaseURL, "/")
	}

	// 从SDK获取的路径，并移除开头可能存在的斜杠
	actionPath := strings.TrimPrefix(a.action, "/")

	// 只有当base_url以'/v1'结尾，且action_path以'v1/'开头时，才进行去重
	if strings.HasSuffix(baseURL, "/v1") && strings.HasPrefix(actionPath, "v1/") {
		actionPath = actionPath[3:] // 从action_path中移除 'v1/'
	}

	finalURL := fmt.Sprintf("%s/%s", baseURL, actionPath)

	// 处理查询参数
	params := make(map[string]string)
	for key, values := range a.c.Request.URL.Query() {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}

	logger.Infof("%s: Rotated to upstream key (masked): %s", a.logPrefix, utils.MaskAPIKeyDefault(upstreamKey))

	return &services.TargetRequest{
		Method:  a.c.Request.Method,
		URL:     finalURL,
		Headers: headers,
		Params:  params,
		Body:    bodyBytes,
	}, nil
}