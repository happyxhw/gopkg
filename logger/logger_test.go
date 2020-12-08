package logger

import (
	"testing"

	"go.uber.org/zap"
)

func TestConsoleLogger(t *testing.T) {
	c := Config{
		Level:       "info",
		FileName:    "",
		EncoderType: "console",
		Caller:      true,
	}

	InitLogger(&c)
	Info("test", zap.String("1", "2"))
}

func TestJsonLogger(t *testing.T) {
	c := Config{
		Level:       "info",
		FileName:    "",
		EncoderType: "json",
		Caller:      true,
	}

	InitLogger(&c)
	Info("test", zap.String("1", "2"))
}

func TestCallerLogger(t *testing.T) {
	c := Config{
		Level:       "info",
		FileName:    "",
		EncoderType: "json",
		Caller:      false,
	}

	InitLogger(&c)
	Info("test", zap.String("1", "2"))
}

func TestFileLogger(t *testing.T) {
	c := Config{
		Level:       "info",
		FileName:    "./test.log",
		EncoderType: "json",
		Caller:      true,
	}

	InitLogger(&c)
	Info("test", zap.String("1", "2"))
	Error("test", zap.String("1", "2"))
}
