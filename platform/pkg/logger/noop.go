package logger

import (
	"context"

	"go.uber.org/zap"
)

// NoopLogger - логгер-заглушка для тестов и инициализации
type NoopLogger struct{}

func (l *NoopLogger) Info(ctx context.Context, msg string, fields ...zap.Field)  {}
func (l *NoopLogger) Error(ctx context.Context, msg string, fields ...zap.Field) {}
func (l *NoopLogger) Debug(ctx context.Context, msg string, fields ...zap.Field) {}
func (l *NoopLogger) Warn(ctx context.Context, msg string, fields ...zap.Field)  {}
func (l *NoopLogger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {}
