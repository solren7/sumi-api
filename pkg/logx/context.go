package logx

import (
	"context"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// ---------- Metadata & context key ----------

// Metadata holds system-level tracing information that is automatically
// injected into every log event's "metadata" dict by the metadataHook.
// Extend this struct to carry additional context fields across the application.
type Metadata struct {
	RequestID string `json:"request_id,omitempty"`
	UserID    string `json:"user_id,omitempty"`
}

// metadataCtxKey is an unexported typed key to avoid collision with other packages.
type metadataCtxKey struct{}

// MetadataKey is the context key used to store/retrieve Metadata.
// Use logx.NewContext() to set it, or set it manually:
//
//	ctx = context.WithValue(ctx, logx.MetadataKey, logx.Metadata{...})
var MetadataKey = metadataCtxKey{}

// NewContext stores Metadata into the given context.
func NewContext(ctx context.Context, md Metadata) context.Context {
	return context.WithValue(ctx, MetadataKey, md)
}

// ---------- metadata hook ----------

// metadataHook auto-extracts Metadata struct and OpenTelemetry trace/span IDs
// from the context attached via event.Ctx(), and writes them as a nested
// "metadata" dict in every log event.
type metadataHook struct{}

func (h metadataHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	ctx := e.GetCtx()
	dict := zerolog.Dict()
	hasData := false

	// 1. Metadata struct from context
	if md, ok := ctx.Value(MetadataKey).(Metadata); ok {
		if md.RequestID != "" {
			dict = dict.Str("request_id", md.RequestID)
			hasData = true
		}
		if md.UserID != "" {
			dict = dict.Str("user_id", md.UserID)
			hasData = true
		}
	}

	// 2. OpenTelemetry trace_id + span_id
	spanCtx := trace.SpanFromContext(ctx).SpanContext()
	if spanCtx.IsValid() {
		dict = dict.Str("trace_id", spanCtx.TraceID().String())
		dict = dict.Str("span_id", spanCtx.SpanID().String())
		hasData = true
	}

	if hasData {
		e.Dict("metadata", dict)
	}
}

// ---------- LogEntry (pooled chainable builder) ----------

// LogEntry is a chainable log builder. Obtain one via WithCtx / WithFields /
// WithError, chain additional methods, then call a terminal method (Info,
// Error, …) which writes the log and returns the entry to the pool.
type LogEntry struct {
	ctx    context.Context
	fields map[string]any
	err    error
}

var entryPool = sync.Pool{
	New: func() any { return &LogEntry{} },
}

func newEntry() *LogEntry {
	e := entryPool.Get().(*LogEntry)
	e.ctx = nil
	e.fields = nil
	e.err = nil
	return e
}

func releaseEntry(e *LogEntry) {
	// Clear the map but keep underlying memory for reuse.
	for k := range e.fields {
		delete(e.fields, k)
	}
	e.ctx = nil
	e.err = nil
	entryPool.Put(e)
}

// ---------- logger instance ----------

var logger zerolog.Logger

func init() {
	logger = newLogger("console", os.Stdout)
}

// SetLogger replaces the package-level logger (e.g. for production JSON output).
func SetLogger(l zerolog.Logger) {
	logger = l.Hook(metadataHook{})
}

func Configure(format string) {
	logger = newLogger(format, os.Stdout)
}

func newLogger(format string, out io.Writer) zerolog.Logger {
	var base zerolog.Logger
	if strings.EqualFold(strings.TrimSpace(format), "json") {
		base = zerolog.New(out).With().Timestamp().Logger()
	} else {
		console := zerolog.ConsoleWriter{Out: out, TimeFormat: time.RFC3339}
		base = zerolog.New(console).With().Timestamp().Logger()
	}
	return base.Hook(metadataHook{})
}
