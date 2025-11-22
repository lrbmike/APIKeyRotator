package cache

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// MemoryClient 内存缓存客户端，用于替代 Redis
type MemoryClient struct {
	data map[string]*cacheItem
	mu   sync.RWMutex
}

type cacheItem struct {
	value      interface{}
	expiration time.Time
}

// InitializeMemory 初始化内存缓存
func InitializeMemory() *MemoryClient {
	client := &MemoryClient{
		data: make(map[string]*cacheItem),
	}

	// 启动后台清理过期项的协程
	go client.cleanupExpired()

	return client
}

// Incr 原子性递增操作，返回递增后的值
func (c *MemoryClient) Incr(ctx context.Context, key string) *IntCmd {
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
			newValue = 1
		}
	}

	c.data[key] = &cacheItem{
		value:      newValue,
		expiration: time.Time{}, // 永不过期
	}

	return &IntCmd{val: newValue, err: nil}
}

// Get 获取键值
func (c *MemoryClient) Get(ctx context.Context, key string) *StringCmd {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.data[key]
	if !exists {
		return &StringCmd{val: "", err: fmt.Errorf("redis: nil")}
	}

	// 检查是否过期
	if !item.expiration.IsZero() && time.Now().After(item.expiration) {
		return &StringCmd{val: "", err: fmt.Errorf("redis: nil")}
	}

	if val, ok := item.value.(string); ok {
		return &StringCmd{val: val, err: nil}
	}

	return &StringCmd{val: fmt.Sprintf("%v", item.value), err: nil}
}

// Set 设置键值
func (c *MemoryClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *StatusCmd {
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

	return &StatusCmd{err: nil}
}

// Del 删除键
func (c *MemoryClient) Del(ctx context.Context, keys ...string) *IntCmd {
	c.mu.Lock()
	defer c.mu.Unlock()

	count := int64(0)
	for _, key := range keys {
		if _, exists := c.data[key]; exists {
			delete(c.data, key)
			count++
		}
	}

	return &IntCmd{val: count, err: nil}
}

// Ping 测试连接
func (c *MemoryClient) Ping(ctx context.Context) *StatusCmd {
	return &StatusCmd{err: nil}
}

// cleanupExpired 定期清理过期的缓存项
func (c *MemoryClient) cleanupExpired() {
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

// IntCmd 整数命令结果
type IntCmd struct {
	val int64
	err error
}

func (cmd *IntCmd) Result() (int64, error) {
	return cmd.val, cmd.err
}

func (cmd *IntCmd) Val() int64 {
	return cmd.val
}

func (cmd *IntCmd) Err() error {
	return cmd.err
}

// StringCmd 字符串命令结果
type StringCmd struct {
	val string
	err error
}

func (cmd *StringCmd) Result() (string, error) {
	return cmd.val, cmd.err
}

func (cmd *StringCmd) Val() string {
	return cmd.val
}

func (cmd *StringCmd) Err() error {
	return cmd.err
}

// StatusCmd 状态命令结果
type StatusCmd struct {
	err error
}

func (cmd *StatusCmd) Result() (string, error) {
	if cmd.err != nil {
		return "", cmd.err
	}
	return "OK", nil
}

func (cmd *StatusCmd) Err() error {
	return cmd.err
}