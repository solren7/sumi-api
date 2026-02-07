package logx

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

// PgxTracer 实现 pgx.QueryTracer 接口
type PgxTracer struct{}

func NewPgxTracer() *PgxTracer {
	return &PgxTracer{}
}

// 定义一个私有的 context key，防止冲突
type pgxTraceKey struct{}

// 定义一个结构体，用于在 Start 和 End 之间传递数据
type pgxTraceInfo struct {
	StartTime time.Time
	SQL       string
	Args      []any
}

// TraceQueryStart 在 SQL 执行前调用
func (t *PgxTracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	// 1. 在这里捕获 SQL 和 参数
	info := &pgxTraceInfo{
		StartTime: time.Now(),
		SQL:       data.SQL,  // <--- SQL 在这里获取
		Args:      data.Args, // <--- 参数也在这里
	}

	// 2. 将 info 放入 Context 返回
	// pgx 会把这个 Context 传递给 TraceQueryEnd
	return context.WithValue(ctx, pgxTraceKey{}, info)
}

// TraceQueryEnd 在 SQL 执行后调用
func (t *PgxTracer) TraceQueryEnd(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryEndData) {
	// 1. 从 Context 中取出 info
	info, ok := ctx.Value(pgxTraceKey{}).(*pgxTraceInfo)
	if !ok {
		return // 如果没有 info，说明 Start 没执行或被截断，直接跳过
	}

	// 2. 整理额外信息
	extra := map[string]any{
		"rows_affected": data.CommandTag.RowsAffected(),
	}

	if data.Err != nil {
		extra["error"] = data.Err.Error()
	}

	// 可选：如果你想打印参数，可以加进去 (注意敏感数据脱敏)
	// if len(info.Args) > 0 {
	//    extra["args"] = info.Args
	// }

	// 3. 调用 logx 记录 Span
	// 使用 info.SQL 和 info.StartTime
	AddSpan(ctx, "pgx", info.SQL, info.StartTime, extra)
}
