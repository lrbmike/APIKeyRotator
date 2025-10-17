package redis

import (
	"context"
	"fmt"

	"api-key-rotator/backend/internal/config"

	"github.com/go-redis/redis/v8"
)

// Initialize 初始化Redis连接
func Initialize(cfg *config.Config) (*redis.Client, error) {
	var rdb *redis.Client

	if cfg.RedisPassword != "" {
		rdb = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
			Password: cfg.RedisPassword,
			DB:       0,
		})
	} else {
		rdb = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
			Password: "",
			DB:       0,
		})
	}

	// 测试连接
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return rdb, nil
}