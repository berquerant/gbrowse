package ctxlog

import (
	"context"

	"go.uber.org/zap"
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

// Logger is a nil safe wrapper of zap.Logger.
type Logger interface {
	Error(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	Sync() error
}

// New returns a new Logger.
func New(instance *zap.Logger) Logger {
	return &logger{
		l: instance,
	}
}

type logger struct {
	l *zap.Logger
}

func (l *logger) Sync() error {
	if l.isNil() {
		return nil
	}
	return l.l.Sync()
}

func (l *logger) Info(msg string, fields ...zap.Field) {
	if !l.isNil() {
		l.l.Info(msg, fields...)
	}
}

func (l *logger) Error(msg string, fields ...zap.Field) {
	if !l.isNil() {
		l.l.Error(msg, fields...)
	}
}

func (l *logger) Debug(msg string, fields ...zap.Field) {
	if !l.isNil() {
		l.l.Debug(msg, fields...)
	}
}

func (l *logger) isNil() bool {
	return l == nil || l.l == nil
}
