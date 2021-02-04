package rabbitmq

import (
	"testing"
	"time"

	"github.com/streadway/amqp"
)

const (
	url    = "amqp://happyxhw:808258@localhost:5672/"
	exName = "gopkg.product"
	exType = amqp.ExchangeDirect

	consumerName = "gopkg_recv"
	queueName    = "product"
	key          = "gopkg"
)

func TestProducer_Publish(t *testing.T) {
	p, err := NewProducer(url, exName, exType)
	if err != nil {
		t.Fatal(err)
	}
	msg := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte("hello world"),
	}
	for {
		_, err = p.Publish(msg, key)
		if err != nil {
			t.Error(err)
		}
		time.Sleep(time.Second)
	}
}
