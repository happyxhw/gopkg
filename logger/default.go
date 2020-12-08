package logger

import (
	"io"
	"log"
	"os"
	"strings"
	"time"

	rotateLogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// encoderConfig 编码控制
var encoderConfig = zapcore.EncoderConfig{
	TimeKey:        "time",
	LevelKey:       "level",
	NameKey:        "logger",
	CallerKey:      "caller",
	MessageKey:     "msg",
	StacktraceKey:  "stacktrace",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.LowercaseLevelEncoder,
	EncodeTime:     timeEncoder,
	EncodeDuration: zapcore.SecondsDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
	EncodeName:     zapcore.FullNameEncoder,
}

// SetUp logger
func SetUp(logLevel zapcore.Level, filename, encoderType string, opts ...zap.Option) *zap.Logger {
	level := zap.NewAtomicLevel()
	level.SetLevel(logLevel)
	var encoder zapcore.Encoder
	if encoderType == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
	var core zapcore.Core
	if filename != "" {
		var cores []zapcore.Core
		writer := Writer(filename)
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(writer), level))
		if logLevel == zap.InfoLevel || logLevel == zap.DebugLevel {
			filename = strings.TrimSuffix(filename, ".log")
			errWriter := Writer(filename + "_err.log")
			cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(errWriter), zap.WarnLevel))
		}
		core = zapcore.NewTee(cores...)
	} else {
		core = zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level)
	}

	return zap.New(core, opts...)
}

func Writer(filename string) io.Writer {
	hook, err := rotateLogs.New(
		filename+".%Y-%m-%d",
		rotateLogs.WithLinkName(filename),
		rotateLogs.WithMaxAge(time.Hour*24*7),
		rotateLogs.WithRotationTime(time.Hour*24),
	)
	if err != nil {
		log.Fatalf("init rotatelogs err: %+v", err)
	}
	return hook
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	layout := "2006-01-02 15:04:05"
	enc.AppendString(t.Format(layout))
}
