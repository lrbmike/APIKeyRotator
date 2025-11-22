package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"api-key-rotator/backend/internal/adapters"
	"api-key-rotator/backend/internal/cache"
	"api-key-rotator/backend/internal/config"
	"api-key-rotator/backend/internal/logger"
	"api-key-rotator/backend/internal/models"
	"api-key-rotator/backend/internal/services"
	"api-key-rotator/backend/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LLMProxyHandler LLM代理处理器
type LLMProxyHandler struct {
	cfg         *config.Config
	db          *gorm.DB
	cacheClient cache.CacheInterface
}

// NewLLMProxyHandler 创建LLM代理处理器实例
func NewLLMProxyHandler(cfg *config.Config, db *gorm.DB, cacheClient cache.CacheInterface) *LLMProxyHandler {
	return &LLMProxyHandler{
		cfg:         cfg,
		db:          db,
		cacheClient: cacheClient,
	}
}

// HandleLLMProxy 处理LLM代理请求
func (h *LLMProxyHandler) HandleLLMProxy(c *gin.Context) {
	slug := c.Param("slug")
	action := strings.TrimPrefix(c.Param("action"), "/")
	
	if err := services.ValidateSlug(slug); err != nil {
		logger.Warningf("Bad Request for LLM slug '%s': %v", slug, err)
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	targetRequest, err := h.prepareLLMRequest(c, slug, action)
	if err != nil {
		logger.Warningf("Bad Request for LLM slug '%s': %v", slug, err)
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	// 转发请求
	if err := h.forwardLLMRequest(c, targetRequest); err != nil {
		logger.Errorf("An unexpected error occurred in LlmApiProxyHandler for slug '%s': %v", slug, err)
		c.JSON(http.StatusBadGateway, gin.H{"detail": "Bad Gateway"})
		return
	}
}

// prepareLLMRequest 准备LLM代理请求
func (h *LLMProxyHandler) prepareLLMRequest(c *gin.Context, slug, action string) (*services.TargetRequest, error) {
	// 1. 加载基础配置以获取 api_format，这是选择适配器的关键
	var proxyConfig models.ProxyConfig
	if err := h.db.Preload("APIKeys").Where("slug = ? AND is_active = ? AND config_type = ?", slug, true, "llm").First(&proxyConfig).Error; err != nil {
		return nil, fmt.Errorf("LLM service configuration with slug '%s' not found or inactive", slug)
	}

	// 2. 根据 api_format 选择适配器
	var adapter adapters.LLMAdapter
	apiFormat := "openai_compatible" // 默认值
	if proxyConfig.APIFormat != nil {
		apiFormat = *proxyConfig.APIFormat
	}

	switch apiFormat {
	case "openai_compatible":
		adapter = adapters.NewOpenAIAdapter(h.cfg, h.db, h.cacheClient, c, &proxyConfig, action)
	case "gemini_native":
		adapter = adapters.NewGeminiAdapter(h.cfg, h.db, h.cacheClient, c, &proxyConfig, action)
	case "anthropic_native":
		adapter = adapters.NewAnthropicAdapter(h.cfg, h.db, h.cacheClient, c, &proxyConfig, action)
	default:
		logger.Errorf("No adapter found for API format '%s'", apiFormat)
		return nil, fmt.Errorf("unsupported API format '%s' for LLM service '%s'", apiFormat, slug)
	}

	// 3. 实例化适配器并移交全部控制权
	logger.Infof("LLM Proxy Handler for '%s': Dispatching to adapter: %s", slug, apiFormat)
	return adapter.ProcessRequest()
}

// forwardLLMRequest 转发LLM请求到目标服务器
func (h *LLMProxyHandler) forwardLLMRequest(c *gin.Context, target *services.TargetRequest) error {
	// 构建目标URL
	targetURL, err := url.Parse(target.URL)
	if err != nil {
		return fmt.Errorf("invalid target URL: %w", err)
	}

	// 添加查询参数
	if len(target.Params) > 0 {
		query := targetURL.Query()
		for key, value := range target.Params {
			query.Set(key, value)
		}
		targetURL.RawQuery = query.Encode()
	}

	// 创建HTTP请求
	req, err := http.NewRequest(target.Method, targetURL.String(), bytes.NewReader(target.Body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	for key, value := range target.Headers {
		req.Header.Set(key, value)
	}

	logger.Infof("Forwarding request to: %s %s", req.Method, req.URL.String())

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	logger.Infof("Received response from target with status code: %d", resp.StatusCode)

	// 过滤响应头
	filteredHeaders := utils.FilterResponseHeaders(resp.Header)
	
	// 设置响应头
	for key, value := range filteredHeaders {
		c.Header(key, value)
	}

	// 设置状态码
	c.Status(resp.StatusCode)

	// 检查是否为流式响应
	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "text/event-stream") {
		// 流式响应
		c.Stream(func(w io.Writer) bool {
			buffer := make([]byte, 1024)
			n, err := resp.Body.Read(buffer)
			if err != nil {
				if err != io.EOF {
					logger.Errorf("Error reading stream: %v", err)
				}
				return false
			}
			_, err = w.Write(buffer[:n])
			return err == nil
		})
	} else {
		// 普通响应
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}
		c.Data(resp.StatusCode, contentType, body)
	}

	return nil
}