package dbgo

import (
	"testing"

	"github.com/happyxhw/gopkg/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestMysql(t *testing.T) {
	l := logger.SetUp(zapcore.InfoLevel, "", "console")
	db, err := NewPostgresDb(&Config{
		User:         "happyxhw",
		Password:     "808258XXxx",
		Host:         "127.0.0.1",
		Port:         "5432",
		DB:           "stravadb",
		MaxIdleConns: 10,
		MaxOpenConns: 10,
		Logger:       l.WithOptions(zap.AddCallerSkip(3), zap.AddCaller()),
		Level:        "info",
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	db.Exec("select idx from t_test")
}
