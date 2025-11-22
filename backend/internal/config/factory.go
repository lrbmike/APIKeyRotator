package config

import (
	"api-key-rotator/backend/internal/infrastructure/cache"
	"api-key-rotator/backend/internal/infrastructure/cache/memory"
	"api-key-rotator/backend/internal/infrastructure/cache/redis"
	"api-key-rotator/backend/internal/infrastructure/database"
	"api-key-rotator/backend/internal/infrastructure/database/mysql"
	"api-key-rotator/backend/internal/infrastructure/database/sqlite"
)

// InfrastructureFactory 基础设施工厂
type InfrastructureFactory struct {
	config *Config
}

// NewInfrastructureFactory 创建基础设施工厂
func NewInfrastructureFactory(cfg *Config) *InfrastructureFactory {
	return &InfrastructureFactory{
		config: cfg,
	}
}

// CreateDatabaseRepository 根据配置创建数据库仓库
func (f *InfrastructureFactory) CreateDatabaseRepository() (database.Repository, error) {
	switch f.config.DBType {
	case "sqlite":
		manager := sqlite.NewSQLiteManager(f.config.DatabasePath)
		return manager.Initialize()
	case "mysql":
		manager := mysql.NewMySQLManager()
		return manager.Initialize()
	default:
		// 默认使用SQLite
		manager := sqlite.NewSQLiteManager(f.config.DatabasePath)
		return manager.Initialize()
	}
}

// CreateCacheInterface 根据配置创建缓存接口
func (f *InfrastructureFactory) CreateCacheInterface() (cache.CacheInterface, error) {
	switch f.config.CacheType {
	case "memory":
		manager := memory.NewMemoryManager()
		return manager.Initialize()
	case "redis":
		manager := redis.NewRedisManager()
		return manager.Initialize()
	default:
		// 默认使用内存缓存
		manager := memory.NewMemoryManager()
		return manager.Initialize()
	}
}