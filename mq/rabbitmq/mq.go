package rabbitmq

import (
	"errors"
)

var (
	ReturnErr  = errors.New("return err")
	NAckErr    = errors.New("nack")
	PublishErr = errors.New("publish err")
	CloseErr   = errors.New("connection closed")
)
