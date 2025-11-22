package database

import (
	"api-key-rotator/backend/internal/models"
	"gorm.io/gorm"
)

// Repository 数据库仓库接口
type Repository interface {
	// 基础操作
	GetDB() *gorm.DB

	// 代理配置管理
	CreateProxyConfig(config *models.ProxyConfig) error
	GetProxyConfigByID(id uint) (*models.ProxyConfig, error)
	GetProxyConfigBySlug(slug string) (*models.ProxyConfig, error)
	UpdateProxyConfig(config *models.ProxyConfig) error
	DeleteProxyConfig(id uint) error
	ListProxyConfigs() ([]*models.ProxyConfig, error)

	// API密钥管理
	CreateAPIKey(key *models.APIKey) error
	GetAPIKeyByID(id uint) (*models.APIKey, error)
	GetAPIKeysByServiceSlug(services string) ([]*models.APIKey, error)
	UpdateAPIKey(key *models.APIKey) error
	DeleteAPIKey(id uint) error
	ListAPIKeys() ([]*models.APIKey, error)

	// 统计和查询
	GetAPIKeyCountByService(serviceSlug string) (int64, error)

	// 数据库迁移和重置
	Migrate() error
	Reset() error
}

// Manager 数据库管理器接口
type Manager interface {
	Initialize() (Repository, error)
	Close() error
	Ping() error
}