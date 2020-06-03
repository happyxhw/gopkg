package logger

import (
	"testing"

	"go.uber.org/zap"
)

func TestInfoDev(t *testing.T) {
	_ = SetupLogger(&Config{
		Level: "info",
		Dev:   true,
	})

	Info("hello", zap.String("h", "w"))
}

func TestInfoProd(t *testing.T) {
	_ = SetupLogger(&Config{
		Level: "info",
		Dev:   false,
	})

	Info("hello", zap.String("h", "w"))
}

func TestErrorDev(t *testing.T) {
	_ = SetupLogger(&Config{
		Level: "info",
		Dev:   true,
	})

	Error("hello", zap.String("h", "w"))
}

func TestErrorProd(t *testing.T) {
	_ = SetupLogger(&Config{
		Level: "info",
		Dev:   false,
	})

	Error("hello", zap.String("h", "w"))
}
