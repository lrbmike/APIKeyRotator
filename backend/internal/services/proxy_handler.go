package services

import (
	"context"
	"fmt"

	"api-key-rotator/backend/internal/cache"
	"api-key-rotator/backend/internal/config"
	"api-key-rotator/backend/internal/logger"
	"api-key-rotator/backend/internal/models"
	"api-key-rotator/backend/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TargetRequest 封装准备好的、即将被转发的请求信息
type TargetRequest struct {
	Method  string
	URL     string
	Headers map[string]string
	Params  map[string]string
	Body    []byte
}

// BaseProxyHandler 代理处理器的抽象基类
type BaseProxyHandler struct {
	cfg         *config.Config
	db          *gorm.DB
	cacheClient cache.CacheInterface
	C           *gin.Context // 导出字段
	Slug        string       // 导出字段
	action      string
	logPrefix   string
}

// NewBaseProxyHandler 创建基础代理处理器
func NewBaseProxyHandler(cfg *config.Config, db *gorm.DB, cacheClient cache.CacheInterface, c *gin.Context, slug, action string) *BaseProxyHandler {
	return &BaseProxyHandler{
		cfg:         cfg,
		db:          db,
		cacheClient: cacheClient,
		C:           c,
		Slug:        slug,
		action:      action,
		logPrefix:   fmt.Sprintf("Proxy Handler for '%s'", slug),
	}
}

// RotateAPIKey 从与给定配置关联的密钥池中轮询一个API Key
func (h *BaseProxyHandler) RotateAPIKey(serviceConfig *models.ProxyConfig) (string, error) {
	// 获取活跃的密钥
	var activeKeys []models.APIKey
	for _, key := range serviceConfig.APIKeys {
		if key.IsActive {
			activeKeys = append(activeKeys, key)
		}
	}

	if len(activeKeys) == 0 {
		logger.Errorf("%s: Service '%s' has no active API keys.", h.logPrefix, serviceConfig.Name)
		return "", fmt.Errorf("no active API keys for this service")
	}

	// 使用缓存原子性递增来实现轮询
	ctx := context.Background()
	keyIndexKey := fmt.Sprintf("proxy_config:%d:key_index", serviceConfig.ID)
	keyIndex, err := h.cacheClient.Incr(ctx, keyIndexKey).Result()
	if err != nil {
		logger.Errorf("%s: Failed to increment key index in cache: %v", h.logPrefix, err)
		return "", fmt.Errorf("failed to rotate API key")
	}

	// 计算实际索引
	actualIndex := int(keyIndex-1) % len(activeKeys)
	selectedKey := activeKeys[actualIndex].KeyValue

	logger.Infof("%s: Selected API key (masked): %s", h.logPrefix, utils.MaskAPIKeyDefault(selectedKey))
	return selectedKey, nil
}

// ValidateSlug 验证slug格式
func ValidateSlug(slug string) error {
	if slug == "" {
		return fmt.Errorf("slug cannot be empty")
	}
	// 可以添加更多验证规则
	return nil
}