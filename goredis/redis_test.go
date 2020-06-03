package goredis

import (
	"fmt"
	"testing"
	"time"
)

func TestNewRedis(t *testing.T) {
	client, err := NewRedis(&Config{
		Host: "127.0.0.1:6379",
	})
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	_, _ = client.Set("test", 1, time.Second*5).Result()

	res, err := client.Get("test").Result()
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if res != "1" {
		t.Fail()
	}
	x, _ := client.Get("test1").Result()
	fmt.Println(x == "")
}
