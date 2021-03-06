package rabbitmq

import (
	"time"

	"github.com/happyxhw/gopkg/logger"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type Consumer struct {
	url string

	conn    *amqp.Connection
	channel *amqp.Channel

	name         string
	queueName    string
	exchangeName string
	key          string

	closeCh chan *amqp.Error
	doneCh  chan struct{}

	retry             int
	initRetryInterval int
}

func NewConsumer(url, name, queueName, exName, key string) (*Consumer, error) {
	c := Consumer{
		url:          url,
		name:         name,
		queueName:    queueName,
		exchangeName: exName,
		key:          key,

		closeCh: make(chan *amqp.Error, 1),
		doneCh:  make(chan struct{}),

		retry:             3,
		initRetryInterval: 3,
	}
	var err error
	c.conn, err = amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	c.channel, err = c.conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := c.channel.QueueDeclare(
		c.queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}
	err = c.channel.QueueBind(
		q.Name,
		c.key,
		c.exchangeName,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	c.channel.NotifyClose(c.closeCh)
	return &c, nil
}

func (c *Consumer) Start(handler func(delivery amqp.Delivery)) error {
	msgs, err := c.channel.Consume(
		c.queueName,
		c.name, // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return err
	}
	for {
		select {
		case m := <-msgs:
			handler(m)
		case <-c.closeCh:
			logger.Error("channel close")
			return CloseErr
		case <-c.doneCh:
			return nil
		}
	}
}

func (c *Consumer) Close() {
	c.doneCh <- struct{}{}
	_ = c.channel.Close()
	_ = c.conn.Close()
}

func (c *Consumer) Reconnect() error {
	var err error
	for i := 0; i < c.retry; i++ {
		logger.Info("reconnecting", zap.Int("retry", i+1))
		time.Sleep(time.Second * time.Duration(c.initRetryInterval*(i+1)))
		_ = c.channel.Close()
		_ = c.conn.Close()
		c.conn, err = amqp.Dial(c.url)
		if err != nil {
			continue
		}
		c.channel, err = c.conn.Channel()
		c.closeCh = make(chan *amqp.Error)
		c.channel.NotifyClose(c.closeCh)
		logger.Info("reconnect successful")
		break
	}
	return err
}
