package adapters

import (
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

// GeminiAdapter 适配器，用于处理Google Gemini原生API格式
type GeminiAdapter struct {
	*BaseLLMAdapter
}

// NewGeminiAdapter 创建Gemini适配器实例
func NewGeminiAdapter(cfg *config.Config, db *gorm.DB, redisClient *redis.Client,
	c *gin.Context, proxyConfig *models.ProxyConfig, action string) *GeminiAdapter {
	return &GeminiAdapter{
		BaseLLMAdapter: NewBaseLLMAdapter(cfg, db, redisClient, c, proxyConfig, action),
	}
}

// ProcessRequest 处理Gemini格式的请求
func (a *GeminiAdapter) ProcessRequest() (*services.TargetRequest, error) {
	// 1. 代理访问认证 (劫持 'x-goog-api-key' Header)
	proxyKey := a.c.GetHeader("x-goog-api-key")
	
	validKeys := a.cfg.GetGlobalProxyKeys()
	isValidKey := false
	for _, key := range validKeys {
		if proxyKey == key {
			isValidKey = true
			break
		}
	}
	if !isValidKey {
		return nil, fmt.Errorf("invalid Proxy Key. Provide it via the 'key' URL query parameter")
	}
	
	// 2. 轮询上游密钥
	upstreamKey, err := a.RotateUpstreamKey()
	if err != nil {
		return nil, err
	}

	// 3. 构建目标请求 (偷梁换柱)
	headers := utils.FilterRequestHeaders(a.c.Request.Header, []string{"x-goog-api-key"})

	// 注入真实的Gemini Key
	headers["x-goog-api-key"] = upstreamKey

	// URL拼接方式也不同
	baseURL := ""
	if a.proxyConfig.TargetBaseURL != nil {
		baseURL = strings.TrimSuffix(*a.proxyConfig.TargetBaseURL, "/")
	}
	finalURL := fmt.Sprintf("%s/%s", baseURL, a.action)

	// 处理查询参数
	params := make(map[string]string)
	for key, values := range a.c.Request.URL.Query() {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}

	// 读取请求体
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