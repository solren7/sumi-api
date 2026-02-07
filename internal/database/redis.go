package database

import (
	"context"
	"fmt"

	"fiber/config" // 替换为你实际的 config 包路径

	"github.com/redis/go-redis/v9"
)

func NewRedis(ctx context.Context, cfg *config.Config) (redis.UniversalClient, error) {
	// 1. 构建通用配置项
	opts := &redis.UniversalOptions{
		Addrs:         []string{cfg.RedisConfig.Addr},
		DialTimeout:   cfg.RedisConfig.DialTimeout,
		ReadTimeout:   cfg.RedisConfig.ReadTimeout,
		WriteTimeout:  cfg.RedisConfig.WriteTimeout,
		PoolSize:      cfg.RedisConfig.PoolSize,
		MinIdleConns:  cfg.RedisConfig.MinIdleConns,
		IsClusterMode: cfg.RedisConfig.IsClusterMode,
	}

	// 创建客户端
	rdb := redis.NewUniversalClient(opts)

	// 3. 立即 Ping 一下，确保连接成功 (Fail Fast)
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return rdb, nil
}
