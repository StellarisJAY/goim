package log

import (
	"github.com/stellarisJAY/goim/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	globalLogger *zap.Logger
)

func init() {
	switch config.Config.Environment {
	case config.DevEnv:
		globalLogger, _ = zap.NewDevelopment()
	case config.ProductEnv:
		globalLogger, _ = zap.NewProduction()
	case config.TestEnv:
		globalLogger = zap.NewExample()
	default:
		panic("unknown logging environment")
	}
}

func Info(message string, fields ...zapcore.Field) {
	globalLogger.Info(message, fields...)
}

func Warn(message string, fields ...zapcore.Field) {
	globalLogger.Warn(message, fields...)
}

func Error(message string, fields ...zapcore.Field) {
	globalLogger.Error(message, fields...)
}

func Debug(message string, fields ...zapcore.Field) {
	globalLogger.Debug(message, fields...)
}
