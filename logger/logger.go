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

// Encoder type
type Encoder int8

const (
	// ConsoleEncoder console encode type
	ConsoleEncoder Encoder = iota
	// JSONEncode json encode type
	JSONEncode
)

var logger = SetUp(zapcore.InfoLevel, "", ConsoleEncoder, zap.AddCallerSkip(1), zap.AddCaller())

// InitLogger init default logger
func InitLogger(logLevel zapcore.Level, filepath string, encoderType Encoder, opts ...zap.Option) {
	logger = SetUp(logLevel, filepath, encoderType, opts...)
}

// SetUp logger
func SetUp(logLevel zapcore.Level, filepath string, encoderType Encoder, opts ...zap.Option) *zap.Logger {
	var encoder zapcore.Encoder
	// encoderConfig 编码控制
	encoderConfig := zapcore.EncoderConfig{
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
	level := zap.NewAtomicLevel()
	level.SetLevel(logLevel)
	if encoderType == JSONEncode {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
	var cores []zapcore.Core
	// 输出到文件，按天分割，error级别下的会把err日志单独输出到 _err.log
	if filepath != "" {
		w := writer(filepath)
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(w), level))
		if logLevel == zap.InfoLevel || logLevel == zap.DebugLevel {
			filepath = strings.TrimSuffix(filepath, ".log")
			errW := writer(filepath + "_err.log")
			cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(errW), zap.WarnLevel))
		}
	}
	// 输出到终端
	cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level))
	core := zapcore.NewTee(cores...)
	return zap.New(core, opts...)
}

func writer(filename string) io.Writer {
	hook, err := rotateLogs.New(
		filename+".%Y-%m-%d",
		rotateLogs.WithLinkName(filename),
		//nolint:gomnd
		rotateLogs.WithMaxAge(time.Hour*24*14),
		//nolint:gomnd
		rotateLogs.WithRotationTime(time.Hour*24),
	)
	if err != nil {
		log.Fatalf("init rotatelogs err: %+v", err)
	}
	return hook
}

// func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
// 	layout := "2006-01-02 15:04:05.06"
// 	enc.AppendString(t.Format(layout))
// }

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
