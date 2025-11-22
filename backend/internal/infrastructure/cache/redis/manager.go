package redis

import (
	"api-key-rotator/backend/internal/infrastructure/cache"
	"context"
	"time"
)

// Manager Redis缓存管理器
type Manager struct {
	cache cache.CacheInterface
}

// NewRedisManager 创建Redis缓存管理器实例
func NewRedisManager() *Manager {
	return &Manager{}
}

// Initialize 初始化Redis缓存
func (m *Manager) Initialize() (cache.CacheInterface, error) {
	// 这里应该从配置中获取Redis连接参数
	// 暂时使用默认参数
	addr := "localhost:6379"
	password := ""
	db := 0

	cacheInstance := NewRedisCache(addr, password, db)

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := cacheInstance.Ping(ctx); err != nil {
		return nil, err
	}

	m.cache = cacheInstance
	return cacheInstance, nil
}

// Close 关闭缓存连接
func (m *Manager) Close() error {
	if redisCache, ok := m.cache.(*Cache); ok {
		return redisCache.client.Close()
	}
	return nil
}