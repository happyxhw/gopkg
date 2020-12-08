package config

import (
	"os"

	"github.com/happyxhw/gopkg/logger"
	"github.com/spf13/viper"
)

// InitConfig init config
func InitConfig(configPath string) error {
	if os.Getenv("ENV") == "DEV" {
		viper.SetConfigName("dev")
	} else if os.Getenv("ENV") == "TEST" {
		viper.SetConfigName("test")
	} else if os.Getenv("ENV") == "PROD" {
		viper.SetConfigName("prod")
	} else {
		logger.Fatal("wrong env not in [DEV, TEST, PROD]")
	}
	viper.AddConfigPath(configPath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}
