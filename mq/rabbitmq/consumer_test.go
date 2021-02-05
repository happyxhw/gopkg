package rabbitmq

import (
	"fmt"
	"testing"

	"github.com/streadway/amqp"
)

func TestConsumer_Start(t *testing.T) {
	c, err := NewConsumer(url, consumerName, queueName, exName, key)
	if err != nil {
		t.Fatal(err)
	}
	fn := func(msg amqp.Delivery) {
		fmt.Println(string(msg.Body))
		_ = msg.Ack(false)
	}
	for {
		err := c.Start(fn)
		if err != nil {
			t.Log(err)
			if err == CloseErr {
				if err := c.Reconnect(); err != nil {
					t.Error(err)
					break
				}
			}
		}
	}
}
