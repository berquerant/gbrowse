package ctxlog

import (
	"context"

	"log/slog"
)

type ctxKey string

var ctxKeyValue ctxKey = "ctxKeyValue"

// From finds Logger from context.
func From(ctx context.Context) Logger {
	val := ctx.Value(ctxKeyValue)
	instance, _ := val.(Logger)
	return instance
}

// With embeds Logger in context.
func With(ctx context.Context, instance Logger) context.Context {
	return context.WithValue(ctx, ctxKeyValue, instance)
}

type Logger interface {
	Error(msg string, attrs ...Attr)
	Info(msg string, attrs ...Attr)
	Debug(msg string, attrs ...Attr)
}

// New returns a new Logger.
func New(instance *slog.Logger) Logger {
	return &logger{
		l: instance,
	}
}

type logger struct {
	l *slog.Logger
}

func (l *logger) logAttrs(level slog.Level, msg string, attrs ...Attr) {
	if l.isNil() {
		return
	}
	sAttrs := make([]slog.Attr, len(attrs))
	for i, attr := range attrs {
		sAttrs[i] = slog.Attr(attr)
	}
	l.l.LogAttrs(nil, level, msg, sAttrs...)
}

func (l *logger) Info(msg string, attrs ...Attr) {
	l.logAttrs(slog.LevelInfo, msg, attrs...)
}

func (l *logger) Error(msg string, attrs ...Attr) {
	l.logAttrs(slog.LevelError, msg, attrs...)
}

func (l *logger) Debug(msg string, attrs ...Attr) {
	l.logAttrs(slog.LevelDebug, msg, attrs...)
}

func (l *logger) isNil() bool {
	return l == nil || l.l == nil
}
