package logx

import (
	"context"

	"github.com/rs/zerolog"
)

// ============================================================
// Package-level direct functions (zero allocation, zero overhead)
// ============================================================

func Info(msg string)  { logger.Info().Msg(msg) }
func Debug(msg string) { logger.Debug().Msg(msg) }
func Warn(msg string)  { logger.Warn().Msg(msg) }
func Error(msg string) { logger.Error().Msg(msg) }

func Infof(format string, args ...any)  { logger.Info().Msgf(format, args...) }
func Debugf(format string, args ...any) { logger.Debug().Msgf(format, args...) }
func Warnf(format string, args ...any)  { logger.Warn().Msgf(format, args...) }
func Errorf(format string, args ...any) { logger.Error().Msgf(format, args...) }

// ============================================================
// Package-level builder starters (return pooled *LogEntry)
// ============================================================

// WithCtx creates a new LogEntry carrying the given context.
// The metadataHook will extract request_id and OTel trace/span IDs
// automatically when the terminal method fires.
func WithCtx(ctx context.Context) *LogEntry {
	e := newEntry()
	e.ctx = ctx
	return e
}

// WithFields creates a new LogEntry with the given fields.
func WithFields(fields map[string]any) *LogEntry {
	e := newEntry()
	e.fields = fields
	return e
}

// WithField creates a new LogEntry with a single field.
func WithField(key string, value any) *LogEntry {
	e := newEntry()
	e.fields = map[string]any{key: value}
	return e
}

// WithError creates a new LogEntry carrying the given error.
func WithError(err error) *LogEntry {
	e := newEntry()
	e.err = err
	return e
}

// ============================================================
// LogEntry chainable methods
// ============================================================

// WithCtx attaches a context to the entry.
func (e *LogEntry) WithCtx(ctx context.Context) *LogEntry {
	e.ctx = ctx
	return e
}

// WithError attaches an error to the entry.
func (e *LogEntry) WithError(err error) *LogEntry {
	e.err = err
	return e
}

// WithFields merges fields into the entry.
func (e *LogEntry) WithFields(fields map[string]any) *LogEntry {
	if e.fields == nil {
		e.fields = fields
	} else {
		for k, v := range fields {
			e.fields[k] = v
		}
	}
	return e
}

// WithField adds a single field to the entry.
func (e *LogEntry) WithField(key string, value any) *LogEntry {
	if e.fields == nil {
		e.fields = make(map[string]any, 4)
	}
	e.fields[key] = value
	return e
}

// ============================================================
// LogEntry terminal methods (write log + release to pool)
// ============================================================

func (e *LogEntry) Info(msg string)  { e.write(zerolog.InfoLevel, msg) }
func (e *LogEntry) Debug(msg string) { e.write(zerolog.DebugLevel, msg) }
func (e *LogEntry) Warn(msg string)  { e.write(zerolog.WarnLevel, msg) }
func (e *LogEntry) Error(msg string) { e.write(zerolog.ErrorLevel, msg) }

func (e *LogEntry) Infof(format string, args ...any)  { e.writef(zerolog.InfoLevel, format, args...) }
func (e *LogEntry) Debugf(format string, args ...any) { e.writef(zerolog.DebugLevel, format, args...) }
func (e *LogEntry) Warnf(format string, args ...any)  { e.writef(zerolog.WarnLevel, format, args...) }
func (e *LogEntry) Errorf(format string, args ...any) { e.writef(zerolog.ErrorLevel, format, args...) }

// ============================================================
// Internal write helpers
// ============================================================

// write is the core log-and-release path for fixed messages.
func (e *LogEntry) write(level zerolog.Level, msg string) {
	event := levelToEvent(level)
	if e.ctx != nil {
		event = event.Ctx(e.ctx) // triggers metadataHook
	}
	if e.err != nil {
		event = event.Err(e.err)
	}
	if len(e.fields) > 0 {
		event = event.Fields(e.fields)
	}
	event.Msg(msg)
	releaseEntry(e)
}

// writef is the core log-and-release path for formatted messages.
func (e *LogEntry) writef(level zerolog.Level, format string, args ...any) {
	event := levelToEvent(level)
	if e.ctx != nil {
		event = event.Ctx(e.ctx) // triggers metadataHook
	}
	if e.err != nil {
		event = event.Err(e.err)
	}
	if len(e.fields) > 0 {
		event = event.Fields(e.fields)
	}
	event.Msgf(format, args...)
	releaseEntry(e)
}

// levelToEvent maps a zerolog.Level to an *zerolog.Event on the package logger.
func levelToEvent(level zerolog.Level) *zerolog.Event {
	switch level {
	case zerolog.DebugLevel:
		return logger.Debug()
	case zerolog.WarnLevel:
		return logger.Warn()
	case zerolog.ErrorLevel:
		return logger.Error()
	case zerolog.FatalLevel:
		return logger.Fatal()
	case zerolog.PanicLevel:
		return logger.Panic()
	default:
		return logger.Info()
	}
}
