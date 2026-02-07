package logx

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/redis/go-redis/v9" // 假设使用 v9
)

type RedisTraceHook struct{}

// NewRedisHook 创建钩子实例
func NewRedisHook() *RedisTraceHook {
	return &RedisTraceHook{}
}

// DialHook 建立连接时的钩子 (通常不需要记录 Span，但必须实现)
func (h *RedisTraceHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return next(ctx, network, addr)
	}
}

// ProcessHook 普通命令钩子 (Get, Set 等)
func (h *RedisTraceHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		start := time.Now()

		// 执行真正的 Redis 命令
		err := next(ctx, cmd)

		// --- 自动记录 Span ---
		// 组装命令详情，例如: "get mykey"
		cmdName := cmd.Name()
		cmdString := fmt.Sprintf("%s %v", cmdName, cmd.Args())

		// 如果太长可以截断
		if len(cmdString) > 200 {
			cmdString = cmdString[:200] + "..."
		}

		AddSpan(ctx, "redis", cmdString, start, map[string]any{
			"error": err, // 记录是否出错
		})

		return err
	}
}

// ProcessPipelineHook 管道命令钩子 (Pipeline)
func (h *RedisTraceHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		start := time.Now()

		err := next(ctx, cmds)

		// 简单记录为 pipeline 操作
		cmdSummaries := make([]string, 0, len(cmds))
		for _, cmd := range cmds {
			cmdSummaries = append(cmdSummaries, cmd.Name())
		}

		AddSpan(ctx, "redis-pipeline", strings.Join(cmdSummaries, ", "), start, map[string]any{
			"count": len(cmds),
			"error": err,
		})

		return err
	}
}
