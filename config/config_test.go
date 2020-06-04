package config

import (
	"fmt"
	"os"
	"testing"
)

func TestSetupConfig(t *testing.T) {
	_ = os.Setenv("ENV", "DEV")
	v, err := SetupConfig(".")
	if err != nil {
		panic(err)
	}

	fmt.Println(v.Get("db"))
}

func TestDecode(t *testing.T) {
	type db struct {
		Host int
		User int
	}

	var d db

	_ = os.Setenv("ENV", "PROD")
	v, err := SetupConfig(".")
	if err != nil {
		panic(err)
	}

	Decode(v.Get("db"), &d)
	fmt.Printf("%+v\n", d)
}
