package cache

import (
	"api-key-rotator/backend/internal/config"
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// CacheInterface 缓存接口，兼容Redis和内存缓存
type CacheInterface interface {
	Incr(ctx context.Context, key string) *IntCmd
	Get(ctx context.Context, key string) *StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *StatusCmd
	Del(ctx context.Context, keys ...string) *IntCmd
	Ping(ctx context.Context) *StatusCmd
}

// Initialize 初始化缓存（工厂模式）
func Initialize(cfg *config.Config) (CacheInterface, error) {
	switch cfg.CacheType {
	case "memory":
		return InitializeMemory(), nil
	case "redis":
		return initializeRedis(cfg)
	default:
		return nil, fmt.Errorf("unsupported cache type: %s", cfg.CacheType)
	}
}

// initializeRedis 初始化Redis缓存
func initializeRedis(cfg *config.Config) (CacheInterface, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       0,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisAdapter{client: rdb}, nil
}

// RedisAdapter Redis适配器，实现CacheInterface接口
type RedisAdapter struct {
	client *redis.Client
}

func (r *RedisAdapter) Incr(ctx context.Context, key string) *IntCmd {
	cmd := r.client.Incr(ctx, key)
	return &IntCmd{val: cmd.Val(), err: cmd.Err()}
}

func (r *RedisAdapter) Get(ctx context.Context, key string) *StringCmd {
	cmd := r.client.Get(ctx, key)
	if cmd.Err() != nil {
		return &StringCmd{val: "", err: cmd.Err()}
	}
	return &StringCmd{val: cmd.Val(), err: cmd.Err()}
}

func (r *RedisAdapter) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *StatusCmd {
	cmd := r.client.Set(ctx, key, value, expiration)
	return &StatusCmd{err: cmd.Err()}
}

func (r *RedisAdapter) Del(ctx context.Context, keys ...string) *IntCmd {
	cmd := r.client.Del(ctx, keys...)
	return &IntCmd{val: cmd.Val(), err: cmd.Err()}
}

func (r *RedisAdapter) Ping(ctx context.Context) *StatusCmd {
	cmd := r.client.Ping(ctx)
	return &StatusCmd{err: cmd.Err()}
}