package dbgo

import (
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	// mysql
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DBType db type
type DBType int8

const (
	// MysqlDB mysql
	MysqlDB DBType = iota
	// PgDB postgresql
	PgDB
)

type Config struct {
	User         string
	Password     string
	Host         string
	Port         int
	DB           string
	MaxIdleConns int `mapstructure:"max_idle_conns"`
	MaxOpenConns int `mapstructure:"max_open_conns"`
	MaxLifeTime  int `mapstructure:"max_life_time"`
	Logger       *zap.Logger
	Level        string
}

func NewMysqlDB(dbConfig *Config) (*gorm.DB, error) {
	DB, err := createConnection(dbConfig, MysqlDB)
	return DB, err
}

func NewPostgresDB(dbConfig *Config) (*gorm.DB, error) {
	DB, err := createConnection(dbConfig, PgDB)
	return DB, err
}

func createConnection(dbConfig *Config, dbType DBType) (*gorm.DB, error) {
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
	if dbType == MysqlDB {
		url := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=UTF8&parseTime=true", user, password, host, port, dbName)
		db, err = gorm.Open(mysql.Open(url), &c)
	} else if dbType == PgDB {
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
			sqlDB.SetConnMaxLifetime(time.Duration(dbConfig.MaxLifeTime) * time.Second)
		}
	}
	return db, err
}
