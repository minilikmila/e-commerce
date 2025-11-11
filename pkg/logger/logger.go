package logger

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New creates a zap logger based on the provided environment name ( production or development).
func New(environment string) (*zap.Logger, error) {
	switch strings.ToLower(environment) {
	case "production", "prod":
		return newProductionLogger()
	default:
		return newDevelopmentLogger()
	}
}

func newDevelopmentLogger() (*zap.Logger, error) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return config.Build()
}

func newProductionLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.Encoding = "json"
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	logger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build production logger: %w", err)
	}

	return logger, nil
}

// WithFields returns a new logger with the provided key-value pairs - customize zap logger fields for each request
func WithFields(base *zap.Logger, fields map[string]interface{}) *zap.Logger {
	if base == nil {
		return zap.NewNop()
	}
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return base.With(zapFields...)
}

// Sync flushes any buffered log entries.
func Sync(log *zap.Logger) {
	if log == nil {
		return
	}
	_ = log.Sync()
}

// Must panics if the logger creation fails.
func Must(log *zap.Logger, err error) *zap.Logger {
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize logger: %v\n", err)
		panic(err)
	}
	return log
}
