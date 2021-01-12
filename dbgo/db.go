package dbgo

import (
	"errors"
	"fmt"

	"go.uber.org/zap"

	// mysql
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	User         string
	Password     string
	Host         string
	Port         int
	DB           string
	MaxIdleConns int `mapstructure:"max_idle_conns"`
	MaxOpenConns int `mapstructure:"max_open_conns"`
	Logger       *zap.Logger
	Level        string
}

func NewMysqlDb(dbConfig *Config) (*gorm.DB, error) {
	var err error
	Db, err := createConnection(dbConfig, "mysql")
	if err != nil {
		return nil, err
	}
	return Db, nil
}

func NewPostgresDb(dbConfig *Config) (*gorm.DB, error) {
	var err error
	Db, err := createConnection(dbConfig, "postgres")
	if err != nil {
		return nil, err
	}
	return Db, nil
}

func createConnection(dbConfig *Config, dbType string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	host := dbConfig.Host
	user := dbConfig.User
	dbName := dbConfig.DB
	password := dbConfig.Password
	port := dbConfig.Port
	if port == 0 {
		port = 3306
	}
	if host == "" {
		host = "127.0.0.1"
	}

	c := gorm.Config{}
	if dbConfig.Logger != nil {
		c.Logger = newLogger(dbConfig.Logger, dbConfig.Level)
	}
	if dbType == "mysql" {
		url := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=UTF8&parseTime=true", user, password, host, port, dbName)
		db, err = gorm.Open(mysql.Open(url), &c)
	} else if dbType == "postgres" {
		url := fmt.Sprintf(
			"host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
			host, port, user, dbName, password,
		)
		db, err = gorm.Open(postgres.Open(url), &c)
	} else {
		return nil, errors.New("unknown db type")
	}

	if err == nil {
		if dbConfig.MaxIdleConns != 0 && dbConfig.MaxOpenConns != 0 {
			sqlDB, err2 := db.DB()
			if err2 != nil {
				return nil, err2
			}
			sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
			sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
		}
	}
	return db, err
}
