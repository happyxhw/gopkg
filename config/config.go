package config

import (
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

func SetupConfig(configPath string) (*viper.Viper, error) {
	v, err := initConfig(configPath)
	return v, err
}

func initConfig(configPath string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName("config")
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
}

func Get(v *viper.Viper, key string) interface{} {
	c := v.Get(key)
	if c == nil {
		panic("no config found")
	}
	return c
}
