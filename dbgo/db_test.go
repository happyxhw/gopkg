package dbgo

import (
	"testing"
)

func TestMysql(t *testing.T) {
	db, err := NewMysqlDb(&Config{
		User:         "happyxhw",
		Password:     "808258",
		Host:         "127.0.0.1",
		Port:         "3306",
		DB:           "micro_recommend",
		MaxIdleConns: 10,
		MaxOpenConns: 10,
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	exists := db.HasTable("news")
	if !exists {
		t.Error("table news not exists")
	}
}
