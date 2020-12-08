package logger

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Level       string
	FileName    string `mapstructure:"file_name"`
	EncoderType string `mapstructure:"encoder_type"`
	Caller      bool
}

var logger = SetUp(zapcore.InfoLevel, "", "console", zap.AddCallerSkip(1), zap.AddCaller())

// InitLogger init logger
func InitLogger(c *Config) {
	level := zapcore.InfoLevel
	switch strings.ToLower(c.Level) {
	case "debug":
		level = zapcore.DebugLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	}
	var opts []zap.Option
	if c.Caller {
		opts = append(opts, zap.AddCaller(), zap.AddCallerSkip(1))
	}
	logger = SetUp(level, c.FileName, c.EncoderType, opts...)
}

func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	logger.Panic(msg, fields...)
}

func GetLogger() *zap.Logger {
	return logger
}

func Sync() {
	_ = logger.Sync()
}
