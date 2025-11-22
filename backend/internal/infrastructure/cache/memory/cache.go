package memory

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Cache 内存缓存实现
type Cache struct {
	data map[string]*cacheItem
	mu   sync.RWMutex
}

type cacheItem struct {
	value      interface{}
	expiration time.Time
}

// NewMemoryCache 创建内存缓存实例
func NewMemoryCache() *Cache {
	c := &Cache{
		data: make(map[string]*cacheItem),
	}

	// 启动后台清理过期项的协程
	go c.cleanupExpired()

	return c
}

// Set 设置键值
func (c *Cache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	exp := time.Time{}
	if expiration > 0 {
		exp = time.Now().Add(expiration)
	}

	c.data[key] = &cacheItem{
		value:      value,
		expiration: exp,
	}

	return nil
}

// Get 获取键值
func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.data[key]
	if !exists {
		return "", fmt.Errorf("key not found")
	}

	// 检查是否过期
	if !item.expiration.IsZero() && time.Now().After(item.expiration) {
		delete(c.data, key)
		return "", fmt.Errorf("key not found")
	}

	if val, ok := item.value.(string); ok {
		return val, nil
	}

	return fmt.Sprintf("%v", item.value), nil
}

// Del 删除键
func (c *Cache) Del(ctx context.Context, keys ...string) (int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	count := int64(0)
	for _, key := range keys {
		if _, exists := c.data[key]; exists {
			delete(c.data, key)
			count++
		}
	}

	return count, nil
}

// Exists 检查键是否存在
func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.data[key]
	if !exists {
		return false, nil
	}

	// 检查是否过期
	if !item.expiration.IsZero() && time.Now().After(item.expiration) {
		delete(c.data, key)
		return false, nil
	}

	return true, nil
}

// Incr 原子性递增操作
func (c *Cache) Incr(ctx context.Context, key string) (int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, exists := c.data[key]
	var newValue int64

	if !exists || (exists && !item.expiration.IsZero() && time.Now().After(item.expiration)) {
		// 键不存在或已过期，初始化为 1
		newValue = 1
	} else {
		// 键存在，递增
		if val, ok := item.value.(int64); ok {
			newValue = val + 1
		} else {
			// 尝试转换
			if val, ok := item.value.(string); ok {
				if parsed, err := parseInt64(val); err == nil {
					newValue = parsed + 1
				} else {
					newValue = 1
				}
			} else {
				newValue = 1
			}
		}
	}

	c.data[key] = &cacheItem{
		value:      newValue,
		expiration: time.Time{}, // 永不过期
	}

	return newValue, nil
}

// IncrBy 原子性递增指定值
func (c *Cache) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, exists := c.data[key]
	var newValue int64

	if !exists || (exists && !item.expiration.IsZero() && time.Now().After(item.expiration)) {
		// 键不存在或已过期，初始化为 value
		newValue = value
	} else {
		// 键存在，递增
		if val, ok := item.value.(int64); ok {
			newValue = val + value
		} else {
			// 尝试转换
			if val, ok := item.value.(string); ok {
				if parsed, err := parseInt64(val); err == nil {
					newValue = parsed + value
				} else {
					newValue = value
				}
			} else {
				newValue = value
			}
		}
	}

	c.data[key] = &cacheItem{
		value:      newValue,
		expiration: time.Time{}, // 永不过期
	}

	return newValue, nil
}

// Decr 原子性递减操作
func (c *Cache) Decr(ctx context.Context, key string) (int64, error) {
	return c.IncrBy(ctx, key, -1)
}

// Ping 测试连接
func (c *Cache) Ping(ctx context.Context) error {
	return nil
}

// Info 获取缓存信息
func (c *Cache) Info(ctx context.Context) (map[string]interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return map[string]interface{}{
		"type":   "memory",
		"keys":   len(c.data),
		"driver": "internal-memory-cache",
	}, nil
}

// cleanupExpired 定期清理过期的缓存项
func (c *Cache) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, item := range c.data {
			if !item.expiration.IsZero() && now.After(item.expiration) {
				delete(c.data, key)
			}
		}
		c.mu.Unlock()
	}
}

// parseInt64 辅助函数，将字符串转换为int64
func parseInt64(s string) (int64, error) {
	var result int64
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}