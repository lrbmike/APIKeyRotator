package handlers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"api-key-rotator/backend/internal/adapters"
	"api-key-rotator/backend/internal/config"
	"api-key-rotator/backend/internal/converters"
	"api-key-rotator/backend/internal/infrastructure/cache"
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

	targetRequest, proxyConfig, err := h.prepareLLMRequest(c, slug, action)
	if err != nil {
		logger.Warningf("Bad Request for LLM slug '%s': %v", slug, err)
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	// 转发请求，传入proxyConfig以支持响应格式转换
	if err := h.forwardLLMRequest(c, targetRequest, proxyConfig); err != nil {
		logger.Errorf("An unexpected error occurred in LlmApiProxyHandler for slug '%s': %v", slug, err)
		c.JSON(http.StatusBadGateway, gin.H{"detail": "Bad Gateway"})
		return
	}
}

// prepareLLMRequest 准备LLM代理请求，返回TargetRequest和ProxyConfig
func (h *LLMProxyHandler) prepareLLMRequest(c *gin.Context, slug, action string) (*services.TargetRequest, *models.ProxyConfig, error) {
	// 1. 加载基础配置
	var proxyConfig models.ProxyConfig
	if err := h.db.Preload("APIKeys").Where("slug = ? AND is_active = ? AND config_type = ?", slug, true, "LLM").First(&proxyConfig).Error; err != nil {
		return nil, nil, fmt.Errorf("LLM service configuration with slug '%s' not found or inactive", slug)
	}

	// 2. 获取格式配置
	apiFormat := "openai_compatible"
	if proxyConfig.APIFormat != nil {
		apiFormat = *proxyConfig.APIFormat
	}
	clientFormat := "none"
	if proxyConfig.OutputFormat != nil {
		clientFormat = *proxyConfig.OutputFormat
	}

	// 3. 检查是否需要格式转换
	needConversion := converters.NeedsConversion(clientFormat, apiFormat)

	// 4. 读取请求体
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read request body: %w", err)
	}

	// 5. 如果需要，转换请求格式
	convertedAction := action
	if needConversion {
		logger.Infof("Request format conversion enabled: %s -> %s", clientFormat, apiFormat)

		converter, err := converters.NewConverter(clientFormat, apiFormat)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create converter: %w", err)
		}

		// 转换请求体
		convertedBody, err := converter.ConvertRequest(bodyBytes)
		if err != nil {
			logger.Errorf("Failed to convert request: %v", err)
			return nil, nil, fmt.Errorf("failed to convert request format: %w", err)
		}
		bodyBytes = convertedBody

		// 转换请求路径
		convertedAction = converter.GetTargetPath(action)

		// 如果路径包含 {model} 占位符，从请求体中提取模型名并替换
		if strings.Contains(convertedAction, "{model}") {
			model := extractModelFromBody(bodyBytes)
			if model != "" {
				convertedAction = strings.ReplaceAll(convertedAction, "{model}", model)
			}
		}

		logger.Infof("Converted action path: %s -> %s", action, convertedAction)
	}

	// 6. 将转换后的body放回request
	c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	// 7. 根据 API format 选择适配器 (使用API格式，而非客户端格式)
	var adapter adapters.LLMAdapter
	switch apiFormat {
	case "openai_compatible":
		adapter = adapters.NewOpenAIAdapter(h.cfg, h.db, h.cacheClient, c, &proxyConfig, convertedAction)
	case "gemini_native":
		adapter = adapters.NewGeminiAdapter(h.cfg, h.db, h.cacheClient, c, &proxyConfig, convertedAction)
	case "anthropic_native":
		adapter = adapters.NewAnthropicAdapter(h.cfg, h.db, h.cacheClient, c, &proxyConfig, convertedAction)
	default:
		logger.Errorf("No adapter found for API format '%s'", apiFormat)
		return nil, nil, fmt.Errorf("unsupported API format '%s' for LLM service '%s'", apiFormat, slug)
	}

	// 8. 实例化适配器并处理请求
	logger.Infof("LLM Proxy Handler for '%s': Dispatching to adapter: %s", slug, apiFormat)
	targetRequest, err := adapter.ProcessRequest()
	if err != nil {
		return nil, nil, err
	}

	return targetRequest, &proxyConfig, nil
}

// forwardLLMRequest 转发LLM请求到目标服务器，并应用响应格式转换
func (h *LLMProxyHandler) forwardLLMRequest(c *gin.Context, target *services.TargetRequest, proxyConfig *models.ProxyConfig) error {
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

	// 获取格式配置，创建转换器
	apiFormat := "openai_compatible"
	if proxyConfig.APIFormat != nil {
		apiFormat = *proxyConfig.APIFormat
	}
	clientFormat := "none"
	if proxyConfig.OutputFormat != nil {
		clientFormat = *proxyConfig.OutputFormat
	}

	needConversion := converters.NeedsConversion(clientFormat, apiFormat)
	var converter *converters.Converter
	if needConversion {
		converter, err = converters.NewConverter(apiFormat, clientFormat)
		if err != nil {
			logger.Errorf("Failed to create response converter: %v", err)
			needConversion = false
		} else {
			logger.Infof("Response format conversion enabled: %s -> %s", apiFormat, clientFormat)
		}
	}

	// 过滤响应头
	filteredHeaders := utils.FilterResponseHeaders(resp.Header)

	// 设置响应头
	for key, value := range filteredHeaders {
		c.Header(key, value)
	}

	// 检查是否为流式响应
	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "text/event-stream") {
		// 流式响应处理
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.Status(resp.StatusCode)

		if needConversion && converter != nil {
			// 带转换的流式响应
			return h.forwardStreamWithConversion(c, resp.Body, converter)
		} else {
			// 直接透传流式响应
			return h.forwardStreamDirect(c, resp.Body)
		}
	} else {
		// 普通响应处理
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		if needConversion && converter != nil {
			// 转换响应格式
			convertedBody, err := converter.ConvertResponse(body)
			if err != nil {
				logger.Errorf("Failed to convert response: %v", err)
				// 转换失败时返回原始响应
				c.Data(resp.StatusCode, contentType, body)
				return nil
			}
			c.Data(resp.StatusCode, "application/json", convertedBody)
		} else {
			c.Data(resp.StatusCode, contentType, body)
		}
	}

	return nil
}

// forwardStreamDirect 直接透传流式响应
func (h *LLMProxyHandler) forwardStreamDirect(c *gin.Context, body io.Reader) error {
	c.Stream(func(w io.Writer) bool {
		buffer := make([]byte, 1024)
		n, err := body.Read(buffer)
		if err != nil {
			if err != io.EOF {
				logger.Errorf("Error reading stream: %v", err)
			}
			return false
		}
		_, err = w.Write(buffer[:n])
		return err == nil
	})
	return nil
}

// forwardStreamWithConversion 带格式转换的流式响应
func (h *LLMProxyHandler) forwardStreamWithConversion(c *gin.Context, body io.Reader, converter *converters.Converter) error {
	scanner := bufio.NewScanner(body)
	// 增加缓冲区大小以处理大的SSE消息
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	c.Stream(func(w io.Writer) bool {
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				logger.Errorf("Error scanning stream: %v", err)
			}
			return false
		}

		line := scanner.Text()

		// 跳过空行
		if strings.TrimSpace(line) == "" {
			w.Write([]byte("\n"))
			return true
		}

		// 处理SSE数据行
		if strings.HasPrefix(line, "data: ") {
			payload := strings.TrimPrefix(line, "data: ")

			// 处理 [DONE] 信号
			if payload == "[DONE]" {
				w.Write([]byte("data: [DONE]\n\n"))
				return true
			}

			// 转换JSON payload
			convertedPayload, err := converter.ConvertStreamChunk([]byte(payload))
			if err != nil {
				logger.Errorf("Failed to convert stream chunk: %v", err)
				// 转换失败时透传原始数据
				w.Write([]byte(line + "\n"))
				return true
			}

			// 如果转换结果为nil，跳过这个chunk
			if convertedPayload == nil {
				return true
			}

			// 写入转换后的数据
			w.Write([]byte("data: "))
			w.Write(convertedPayload)
			w.Write([]byte("\n\n"))
			return true
		}

		// 其他行（如event:, id:等）直接透传
		w.Write([]byte(line + "\n"))
		return true
	})

	return nil
}

// extractModelFromBody extracts the model name from a request body
func extractModelFromBody(body []byte) string {
	var req map[string]interface{}
	if err := json.Unmarshal(body, &req); err != nil {
		return ""
	}

	if model, ok := req["model"].(string); ok {
		return model
	}

	return ""
}
