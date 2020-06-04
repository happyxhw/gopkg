package config

import (
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

func SetupConfig(configPath string) (*viper.Viper, error) {
	v, err := initConfig(configPath)
	return v, err
}

func initConfig(configPath string) (*viper.Viper, error) {
	v := viper.New()
	if os.Getenv("ENV") == "DEV" {
		v.SetConfigName("dev")
	} else if os.Getenv("ENV") == "TEST" {
		v.SetConfigName("test")
	} else if os.Getenv("ENV") == "PROD" {
		v.SetConfigName("prod")
	} else {
		panic("wrong env not in [DEV, TEST, PROD]")
	}
	v.AddConfigPath(configPath)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	return v, nil
}

func Decode(src, dst interface{}) {
	if src == nil {
		panic("no config found")
	}
	err := mapstructure.Decode(src, dst)
	if err != nil {
		panic(err)
	}
	if dst == nil {
		panic("no config result")
	}
}

func Get(v *viper.Viper, key string) interface{} {
	c := v.Get(key)
	if c == nil {
		panic("no config found")
	}
	return c
}
