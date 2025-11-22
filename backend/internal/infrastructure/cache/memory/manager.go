package memory

import (
	"api-key-rotator/backend/internal/infrastructure/cache"
)

// Manager 内存缓存管理器
type Manager struct {
	cache cache.CacheInterface
}

// NewMemoryManager 创建内存缓存管理器实例
func NewMemoryManager() *Manager {
	return &Manager{}
}

// Initialize 初始化内存缓存
func (m *Manager) Initialize() (cache.CacheInterface, error) {
	cacheInstance := NewMemoryCache()
	m.cache = cacheInstance
	return cacheInstance, nil
}

// Close 关闭缓存连接
func (m *Manager) Close() error {
	// 内存缓存不需要关闭连接
	return nil
}