package logx

import (
	"context"
	"sync"
	"time"
)

// Span 代表一次具体的内部操作（SQL, Redis, RPC等）
type Span struct {
	Type      string         // 类型: "sql", "redis", "rpc"
	Cmd       string         // 具体命令: "select * from users", "get key"
	Duration  time.Duration  // 耗时
	Extra     map[string]any // 额外信息: rows_affected, error 等
	StartTime time.Time      // 开始时间（用于排序或计算相对时间）
}

// TraceRecorder 用于在 Context 中收集所有的 Span
type TraceRecorder struct {
	Spans []Span
	mu    sync.Mutex
}

// TraceKey context key
const TraceKey = "trace_recorder"

// NewTraceContext 初始化记录器并注入 Context
func NewTraceContext(ctx context.Context) (context.Context, *TraceRecorder) {
	tr := &TraceRecorder{
		Spans: make([]Span, 0, 8), // 预分配一点空间
	}
	return context.WithValue(ctx, TraceKey, tr), tr
}

// AddSpan 往当前 Context 的记录器中追加一条记录
func AddSpan(ctx context.Context, spanType, cmd string, start time.Time, extra map[string]any) {
	val := ctx.Value(TraceKey)
	if val == nil {
		return // 如果没有开启 Trace，直接忽略，不影响业务
	}

	tr, ok := val.(*TraceRecorder)
	if !ok {
		return
	}

	duration := time.Since(start)

	// 线程安全地追加
	tr.mu.Lock()
	tr.Spans = append(tr.Spans, Span{
		Type:      spanType,
		Cmd:       cmd,
		Duration:  duration,
		Extra:     extra,
		StartTime: start,
	})
	tr.mu.Unlock()
}
