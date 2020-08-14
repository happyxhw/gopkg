package logger

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

type Config struct {
	Level string
	Dev   bool
}

func init() {
	_ = InitLogger(zap.InfoLevel, 1, true)
}

// SetupLogger setup logger
func SetupLogger(config *Config) error {
	level := strings.ToUpper(config.Level)
	logLevel := zap.InfoLevel
	if level == "DEBUG" {
		logLevel = zap.DebugLevel
	} else if level == "WARN" {
		logLevel = zap.WarnLevel
	} else if level == "ERROR" {
		logLevel = zap.ErrorLevel
	}

	err := InitLogger(logLevel, 1, config.Dev)
	return err
}

var encoderConfig = zapcore.EncoderConfig{
	TimeKey:        "time",
	LevelKey:       "level",
	NameKey:        "logger",
	CallerKey:      "caller",
	MessageKey:     "msg",
	StacktraceKey:  "stacktrace",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.LowercaseLevelEncoder,
	EncodeTime:     zapcore.ISO8601TimeEncoder,
	EncodeDuration: zapcore.SecondsDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
	EncodeName:     zapcore.FullNameEncoder,
}

// InitLogger initial the zap logger config
func InitLogger(logLevel zapcore.Level, skip int, dev bool) error {
	var err error

	var config zap.Config
	devConfig := zap.NewDevelopmentConfig()
	devConfig.OutputPaths = []string{"stdout"}
	devConfig.ErrorOutputPaths = []string{"stderr"}

	prodConfig := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	if dev {
		config = devConfig
		config.DisableCaller = false
		config.DisableStacktrace = true
	} else {
		config = prodConfig
		config.DisableCaller = true
		config.DisableStacktrace = true
		config.Sampling = nil
	}
	level := zap.NewAtomicLevel()
	level.SetLevel(logLevel)
	config.Level = level
	config.EncoderConfig = encoderConfig

	logger, err = config.Build(zap.AddCallerSkip(skip))
	return err
}

func Sync() {
	_ = logger.Sync()
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
