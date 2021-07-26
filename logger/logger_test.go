package logger_test

import (
	"errors"
	"testing"

	"github.com/happyxhw/gopkg/logger"

	"go.uber.org/zap"
)

func TestError(t *testing.T) {
	logger.InitLogger(zap.InfoLevel, "", logger.ConsoleEncoder, zap.AddCallerSkip(1), zap.AddCaller())
	logger.Error("hello", zap.String("key", "value"), zap.Error(errors.New("test")))
}
