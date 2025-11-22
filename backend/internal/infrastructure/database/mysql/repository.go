package mysql

import (
	"api-key-rotator/backend/internal/models"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Repository MySQL 数据库仓库实现
type Repository struct {
	db *gorm.DB
}

// NewMySQLRepository 创建MySQL仓库实例
func NewMySQLRepository(databaseURL string) (*Repository, error) {
	db, err := gorm.Open(mysql.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL database: %w", err)
	}

	return &Repository{db: db}, nil
}

// GetDB 返回GORM实例
func (r *Repository) GetDB() *gorm.DB {
	return r.db
}

// CreateProxyConfig 创建代理配置
func (r *Repository) CreateProxyConfig(config *models.ProxyConfig) error {
	return r.db.Create(config).Error
}

// GetProxyConfigByID 根据ID获取代理配置
func (r *Repository) GetProxyConfigByID(id uint) (*models.ProxyConfig, error) {
	var config models.ProxyConfig
	err := r.db.First(&config, id).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// GetProxyConfigBySlug 根据服务标识获取代理配置
func (r *Repository) GetProxyConfigBySlug(slug string) (*models.ProxyConfig, error) {
	var config models.ProxyConfig
	err := r.db.Where("service_slug = ?", slug).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// UpdateProxyConfig 更新代理配置
func (r *Repository) UpdateProxyConfig(config *models.ProxyConfig) error {
	return r.db.Save(config).Error
}

// DeleteProxyConfig 删除代理配置
func (r *Repository) DeleteProxyConfig(id uint) error {
	return r.db.Delete(&models.ProxyConfig{}, id).Error
}

// ListProxyConfigs 列出所有代理配置
func (r *Repository) ListProxyConfigs() ([]*models.ProxyConfig, error) {
	var configs []*models.ProxyConfig
	err := r.db.Find(&configs).Error
	return configs, err
}

// CreateAPIKey 创建API密钥
func (r *Repository) CreateAPIKey(key *models.APIKey) error {
	return r.db.Create(key).Error
}

// GetAPIKeyByID 根据ID获取API密钥
func (r *Repository) GetAPIKeyByID(id uint) (*models.APIKey, error) {
	var key models.APIKey
	err := r.db.First(&key, id).Error
	if err != nil {
		return nil, err
	}
	return &key, nil
}

// GetAPIKeysByServiceSlug 根据服务标识获取API密钥列表
func (r *Repository) GetAPIKeysByServiceSlug(serviceSlug string) ([]*models.APIKey, error) {
	var keys []*models.APIKey
	err := r.db.Where("service_slug = ?", serviceSlug).Find(&keys).Error
	return keys, err
}

// UpdateAPIKey 更新API密钥
func (r *Repository) UpdateAPIKey(key *models.APIKey) error {
	return r.db.Save(key).Error
}

// DeleteAPIKey 删除API密钥
func (r *Repository) DeleteAPIKey(id uint) error {
	return r.db.Delete(&models.APIKey{}, id).Error
}

// ListAPIKeys 列出所有API密钥
func (r *Repository) ListAPIKeys() ([]*models.APIKey, error) {
	var keys []*models.APIKey
	err := r.db.Find(&keys).Error
	return keys, err
}

// GetAPIKeyCountByService 获取指定服务的API密钥数量
func (r *Repository) GetAPIKeyCountByService(serviceSlug string) (int64, error) {
	var count int64
	err := r.db.Model(&models.APIKey{}).Where("service_slug = ?", serviceSlug).Count(&count).Error
	return count, err
}

// Migrate 执行数据库迁移
func (r *Repository) Migrate() error {
	return r.db.AutoMigrate(
		&models.ProxyConfig{},
		&models.APIKey{},
	)
}

// Reset 重置数据库表
func (r *Repository) Reset() error {
	// 按依赖顺序删除表
	tables := []interface{}{
		&models.APIKey{},
		&models.ProxyConfig{},
	}

	for _, table := range tables {
		if err := r.db.Migrator().DropTable(table); err != nil {
			// 忽略表不存在的错误
			continue
		}
	}

	// 重新创建表
	return r.Migrate()
}