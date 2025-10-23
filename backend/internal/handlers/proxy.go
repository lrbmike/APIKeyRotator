package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

// ProxyHandler 通用代理处理器
type ProxyHandler struct {
	cfg         *config.Config
	db          *gorm.DB
	redisClient *redis.Client
}

// NewProxyHandler 创建通用代理处理器实例
func NewProxyHandler(cfg *config.Config, db *gorm.DB, redisClient *redis.Client) *ProxyHandler {
	return &ProxyHandler{
		cfg:         cfg,
		db:          db,
		redisClient: redisClient,
	}
}

// HandleGenericProxy 处理通用代理请求
func (h *ProxyHandler) HandleGenericProxy(c *gin.Context) {
	slug := strings.TrimPrefix(c.Param("slug"), "/")
	
	// 提取服务标识符（第一个路径段）
	parts := strings.SplitN(slug, "/", 2)
	serviceSlug := parts[0]
	
	if err := services.ValidateSlug(serviceSlug); err != nil {
		logger.Warningf("Bad Request for slug '%s': %v", serviceSlug, err)
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	handler := services.NewBaseProxyHandler(h.cfg, h.db, h.redisClient, c, serviceSlug, "")
	
	targetRequest, err := h.prepareGenericRequest(handler)
	if err != nil {
		logger.Warningf("Bad Request for slug '%s': %v", serviceSlug, err)
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}
	
	// 将完整的路径传递给转发函数
	c.Set("fullPath", slug)

	// 转发请求
	if err := h.forwardRequest(c, targetRequest); err != nil {
		logger.Errorf("An unexpected error occurred in GenericApiProxyHandler for slug '%s': %v", serviceSlug, err)
		c.JSON(http.StatusBadGateway, gin.H{"detail": "Bad Gateway"})
		return
	}
}

// prepareGenericRequest 准备通用代理请求
func (h *ProxyHandler) prepareGenericRequest(handler *services.BaseProxyHandler) (*services.TargetRequest, error) {
	// 1. 认证 (只支持Header)
	proxyKeyHeader := handler.C.GetHeader("X-Proxy-Key")
	validKeys := h.cfg.GetGlobalProxyKeys()
	isValidKey := false
	for _, key := range validKeys {
		if proxyKeyHeader == key {
			isValidKey = true
			break
		}
	}
	if !isValidKey {
		return nil, fmt.Errorf("invalid or missing X-Proxy-Key header")
	}

	// 2. 加载配置
	var proxyConfig models.ProxyConfig
	if err := h.db.Preload("APIKeys").Where("slug = ? AND is_active = ? AND config_type = ?", handler.Slug, true, "generic").First(&proxyConfig).Error; err != nil {
		return nil, fmt.Errorf("generic service configuration with slug '%s' not found or inactive", handler.Slug)
	}

	// 3. 方法校验
	if proxyConfig.Method == nil || strings.ToUpper(handler.C.Request.Method) != strings.ToUpper(*proxyConfig.Method) {
		return nil, fmt.Errorf("method Not Allowed. This path only accepts %s, but received %s",
			strings.ToUpper(*proxyConfig.Method), strings.ToUpper(handler.C.Request.Method))
	}

	// 4. 轮询并注入密钥
	apiKey, err := handler.RotateAPIKey(&proxyConfig)
	if err != nil {
		return nil, err
	}

	// 5. 处理请求头
	headers := utils.FilterRequestHeaders(handler.C.Request.Header, []string{"x-proxy-key"})

	// 6. 处理查询参数
	params := make(map[string]string)
	for key, values := range handler.C.Request.URL.Query() {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}

	// 7. 根据配置注入API密钥
	if proxyConfig.APIKeyLocation != nil && proxyConfig.APIKeyName != nil {
		location := strings.ToLower(*proxyConfig.APIKeyLocation)
		keyName := *proxyConfig.APIKeyName
		if location == "header" {
			headers[keyName] = apiKey
		} else if location == "query" {
			params[keyName] = apiKey
		}
	}

	// 8. 读取请求体
	body, err := io.ReadAll(handler.C.Request.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}

	return &services.TargetRequest{
		Method:  handler.C.Request.Method,
		URL:     *proxyConfig.TargetURL,
		Headers: headers,
		Params:  params,
		Body:    body,
	}, nil
}

// forwardRequest 转发请求到目标服务器
func (h *ProxyHandler) forwardRequest(c *gin.Context, target *services.TargetRequest) error {
	// 构建目标URL
	targetURL, err := url.Parse(target.URL)
	if err != nil {
		return fmt.Errorf("invalid target URL: %w", err)
	}

	// 获取完整路径并处理
	fullPath, _ := c.Get("fullPath")
	requestPath := fullPath.(string)
	
	// 提取除了服务标识符之外的路径部分
	parts := strings.SplitN(requestPath, "/", 2)
	if len(parts) > 1 {
		requestPath = parts[1]
	} else {
		requestPath = ""
	}

	// 如果目标URL没有以"/"结尾且请求路径不为空，则添加"/"
	if requestPath != "" {
		if !strings.HasSuffix(targetURL.Path, "/") && !strings.HasPrefix(requestPath, "/") {
			targetURL.Path += "/"
		}
		targetURL.Path += requestPath
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