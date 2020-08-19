package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestSetupConfig(t *testing.T) {
	_ = os.Setenv("ENV", "DEV")
	err := InitConfig(".")
	if err != nil {
		panic(err)
	}

	fmt.Println(viper.Get("db"))
}
