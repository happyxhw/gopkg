package rabbitmq

import (
	"testing"
)

func TestConsumer_Start(t *testing.T) {
	c, err := NewConsumer(url, consumerName, queueName, exName, key)
	if err != nil {
		t.Fatal(err)
	}
	if err := c.Start(); err != nil {
		t.Error(err)
	}
}
