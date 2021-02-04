package rabbitmq

import (
	"errors"

	"github.com/happyxhw/gopkg/logger"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

var (
	ReturnErr  = errors.New("return err")
	NAckErr    = errors.New("nack")
	PublishErr = errors.New("publish err")
)

type Producer struct {
	url string

	conn    *amqp.Connection
	channel *amqp.Channel

	exchangeName string
	exchangeType string

	ackCh    chan uint64
	nAckCh   chan uint64
	returnCh chan amqp.Return
}

func NewProducer(url, exName, exType string) (*Producer, error) {
	p := Producer{
		url:          url,
		exchangeName: exName,
		exchangeType: exType,

		ackCh:    make(chan uint64),
		nAckCh:   make(chan uint64),
		returnCh: make(chan amqp.Return, 1),
	}
	var err error
	p.conn, err = amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	p.channel, err = p.conn.Channel()
	if err != nil {
		return nil, err
	}
	err = p.channel.Confirm(false)
	if err != nil {
		return nil, err
	}
	p.channel.NotifyConfirm(p.ackCh, p.nAckCh)
	p.channel.NotifyReturn(p.returnCh)
	err = p.channel.ExchangeDeclare(
		p.exchangeName,
		p.exchangeType,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (p *Producer) Publish(msg amqp.Publishing, key string) (*amqp.Return, error) {
	var err error
	err = p.channel.Publish(
		p.exchangeName,
		key,
		true,
		false,
		msg,
	)
	if err != nil {
		logger.Error("publish", zap.Error(err))
		return nil, err
	}
	select {
	case r := <-p.returnCh:
		return &r, ReturnErr
	case <-p.ackCh:
		return nil, nil
	case <-p.nAckCh:
		return nil, NAckErr
	}
}

func (p *Producer) Close() {
	_ = p.channel.Close()
	_ = p.conn.Close()
}
