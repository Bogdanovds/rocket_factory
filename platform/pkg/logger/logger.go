package logger

import (
	"context"
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Key string

const (
	traceIDKey Key = "trace_id"
	userIDKey  Key = "user_id"
)

// Глобальный singleton логгер
var (
	globalLogger *Logger
	initOnce     sync.Once
	dynamicLevel zap.AtomicLevel
)

// Logger обёртка над zap.Logger с enrich поддержкой контекста
type Logger struct {
	zapLogger *zap.Logger
}

// Init инициализирует глобальный логгер.
func Init(levelStr string, asJSON bool) error {
	initOnce.Do(func() {
		dynamicLevel = zap.NewAtomicLevelAt(parseLevel(levelStr))

		encoderCfg := buildProductionEncoderConfig()

		var encoder zapcore.Encoder
		if asJSON {
			encoder = zapcore.NewJSONEncoder(encoderCfg)
		} else {
			encoder = zapcore.NewConsoleEncoder(encoderCfg)
		}

		core := zapcore.NewCore(
			encoder,
			zapcore.AddSync(os.Stdout),
			dynamicLevel,
		)

		zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))

		globalLogger = &Logger{
			zapLogger: zapLogger,
		}
	})

	return nil
}

func buildProductionEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
}

// SetLevel динамически меняет уровень логирования
func SetLevel(levelStr string) {
	if dynamicLevel == (zap.AtomicLevel{}) {
		return
	}

	dynamicLevel.SetLevel(parseLevel(levelStr))
}

// GetLogger возвращает глобальный enrich-aware логгер
func GetLogger() *Logger {
	return globalLogger
}

// SetNopLogger устанавливает глобальный логгер в no-op режим.
func SetNopLogger() {
	globalLogger = &Logger{
		zapLogger: zap.NewNop(),
	}
}

// Sync сбрасывает буферы логгера
func Sync() error {
	if globalLogger != nil {
		return globalLogger.zapLogger.Sync()
	}

	return nil
}

// With создает новый enrich-aware логгер с дополнительными полями
func With(fields ...zap.Field) *Logger {
	if globalLogger == nil {
		return &Logger{zapLogger: zap.NewNop()}
	}

	return &Logger{
		zapLogger: globalLogger.zapLogger.With(fields...),
	}
}

// WithContext создает enrich-aware логгер с контекстом
func WithContext(ctx context.Context) *Logger {
	if globalLogger == nil {
		return &Logger{zapLogger: zap.NewNop()}
	}

	return &Logger{
		zapLogger: globalLogger.zapLogger.With(fieldsFromContext(ctx)...),
	}
}

// Debug enrich-aware debug log
func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	if globalLogger != nil {
		globalLogger.Debug(ctx, msg, fields...)
	}
}

// Info enrich-aware info log
func Info(ctx context.Context, msg string, fields ...zap.Field) {
	if globalLogger != nil {
		globalLogger.Info(ctx, msg, fields...)
	}
}

// Warn enrich-aware warn log
func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	if globalLogger != nil {
		globalLogger.Warn(ctx, msg, fields...)
	}
}

// Error enrich-aware error log
func Error(ctx context.Context, msg string, fields ...zap.Field) {
	if globalLogger != nil {
		globalLogger.Error(ctx, msg, fields...)
	}
}

// Fatal enrich-aware fatal log
func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	if globalLogger != nil {
		globalLogger.Fatal(ctx, msg, fields...)
	}
}

// Instance methods для enrich loggers (Logger)

func (l *Logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Debug(msg, allFields...)
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Info(msg, allFields...)
}

func (l *Logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Warn(msg, allFields...)
}

func (l *Logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Error(msg, allFields...)
}

func (l *Logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Fatal(msg, allFields...)
}

// parseLevel конвертирует строковый уровень в zapcore.Level
func parseLevel(levelStr string) zapcore.Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// fieldsFromContext вытаскивает enrich-поля из контекста
func fieldsFromContext(ctx context.Context) []zap.Field {
	fields := make([]zap.Field, 0)

	if traceID, ok := ctx.Value(traceIDKey).(string); ok && traceID != "" {
		fields = append(fields, zap.String(string(traceIDKey), traceID))
	}

	if userID, ok := ctx.Value(userIDKey).(string); ok && userID != "" {
		fields = append(fields, zap.String(string(userIDKey), userID))
	}

	return fields
}
