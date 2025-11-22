package adapters

import (
	"fmt"

	"api-key-rotator/backend/internal/infrastructure/cache"
	"api-key-rotator/backend/internal/config"
	"api-key-rotator/backend/internal/models"
	"api-key-rotator/backend/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LLMAdapter LLM适配器接口
type LLMAdapter interface {
	ProcessRequest() (*services.TargetRequest, error)
}

// BaseLLMAdapter LLM适配器的抽象基类
type BaseLLMAdapter struct {
	cfg         *config.Config
	db          *gorm.DB
	cacheClient cache.CacheInterface
	c           *gin.Context
	proxyConfig *models.ProxyConfig
	action      string
	logPrefix   string
}

// NewBaseLLMAdapter 创建基础LLM适配器
func NewBaseLLMAdapter(cfg *config.Config, db *gorm.DB, cacheClient cache.CacheInterface,
	c *gin.Context, proxyConfig *models.ProxyConfig, action string) *BaseLLMAdapter {
	apiFormat := "unknown"
	if proxyConfig.APIFormat != nil {
		apiFormat = *proxyConfig.APIFormat
	}
	return &BaseLLMAdapter{
		cfg:         cfg,
		db:          db,
		cacheClient: cacheClient,
		c:           c,
		proxyConfig: proxyConfig,
		action:      action,
		logPrefix:   "Adapter (" + apiFormat + " for ID:" + fmt.Sprintf("%d", proxyConfig.ID) + ")",
	}
}

// RotateUpstreamKey 从密钥池中轮询一个真实的上游API Key
func (a *BaseLLMAdapter) RotateUpstreamKey() (string, error) {
	// 直接使用预加载好的ProxyConfig
	handler := services.NewBaseProxyHandler(a.cfg, a.db, a.cacheClient, a.c, a.proxyConfig.Slug, a.action)
	return handler.RotateAPIKey(a.proxyConfig)
}