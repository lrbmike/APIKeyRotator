package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// Cache Redis缓存实现
type Cache struct {
	client *redis.Client
}

// NewRedisCache 创建Redis缓存实例
func NewRedisCache(addr string, password string, db int) *Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &Cache{client: rdb}
}

// Set 设置键值
func (c *Cache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

// Get 获取键值
func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	result, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

// Del 删除键
func (c *Cache) Del(ctx context.Context, keys ...string) (int64, error) {
	result, err := c.client.Del(ctx, keys...).Result()
	if err != nil {
		return 0, err
	}
	return result, nil
}

// Exists 检查键是否存在
func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	result, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// Incr 原子性递增操作
func (c *Cache) Incr(ctx context.Context, key string) (int64, error) {
	result, err := c.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return result, nil
}

// IncrBy 原子性递增指定值
func (c *Cache) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	result, err := c.client.IncrBy(ctx, key, value).Result()
	if err != nil {
		return 0, err
	}
	return result, nil
}

// Decr 原子性递减操作
func (c *Cache) Decr(ctx context.Context, key string) (int64, error) {
	result, err := c.client.Decr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return result, nil
}

// Ping 测试连接
func (c *Cache) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

// Info 获取缓存信息
func (c *Cache) Info(ctx context.Context) (map[string]interface{}, error) {
	info, err := c.client.Info(ctx).Result()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"type":   "redis",
		"info":   info,
		"driver": "go-redis",
	}, nil
}