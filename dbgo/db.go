package dbgo

import (
	"errors"
	"fmt"
	"github.com/happyxhw/gopkg/logger"

	// mysql
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/wantedly/gorm-zap"
)

type Config struct {
	User         string
	Password     string
	Host         string
	Port         string
	DB           string
	MaxIdleConns int `mapstructure:"max_idle_conns"`
	MaxOpenConns int `mapstructure:"max_open_conns"`
	Log          bool
}

func NewMysqlDb(dbConfig *Config) (*gorm.DB, error) {
	var err error
	Db, err := createConnection(dbConfig, "mysql")
	if err != nil {
		return nil, err
	}
	Db.LogMode(dbConfig.Log)
	Db.SetLogger(gormzap.New(logger.GetLogger()))
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

	if dbType == "mysql" {
		url := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=UTF8&parseTime=true", user, password, host, port, dbName)
		db, err = gorm.Open("mysql", url)
	} else if dbType == "postgres" {
		url := fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
			host, port, user, dbName, password,
		)
		db, err = gorm.Open("postgres", url)
	} else {
		return nil, errors.New("unknown db type")
	}

	if err == nil {
		if dbConfig.MaxIdleConns != 0 && dbConfig.MaxOpenConns != 0 {
			db.DB().SetMaxIdleConns(dbConfig.MaxIdleConns)
			db.DB().SetMaxOpenConns(dbConfig.MaxOpenConns)
		}
	}

	return db, err
}
