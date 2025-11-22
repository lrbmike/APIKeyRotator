package cache

import (
	"context"
	"time"
)

// CacheInterface 缓存接口
type CacheInterface interface {
	// 基础操作
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, keys ...string) (int64, error)
	Exists(ctx context.Context, key string) (bool, error)

	// 计数器操作
	Incr(ctx context.Context, key string) (int64, error)
	IncrBy(ctx context.Context, key string, value int64) (int64, error)
	Decr(ctx context.Context, key string) (int64, error)

	// 连接测试
	Ping(ctx context.Context) error

	// 健康检查和统计
	Info(ctx context.Context) (map[string]interface{}, error)
}

// Manager 缓存管理器接口
type Manager interface {
	Initialize() (CacheInterface, error)
	Close() error
}