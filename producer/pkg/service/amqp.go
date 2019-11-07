package service

import (
	"db-transfer/producer/flags"
	"fmt"
	"github.com/streadway/amqp"
)

const (
	exchangeName = "transfer"
)

type MQ struct {
	channel    *amqp.Channel
	Connection *amqp.Connection
}

func NewMQ(conf *flags.Config) (*MQ, error) {
	var err error
	var connection *amqp.Connection
	var channel *amqp.Channel

	URI := fmt.Sprintf("amqp://%s:%s@%s:%d/", conf.AmqpUser, conf.AmqpPassword, conf.AmqpHost, conf.AmqpPort)
	if connection, err = amqp.Dial(URI); err != nil {
		return nil, fmt.Errorf("amqp connection failed err: %s", err)
	}

	if channel, err = connection.Channel(); err != nil {
		return nil, fmt.Errorf("amqp connection failed err: %s", err)
	}

	if err = channel.ExchangeDeclare(
		exchangeName,
		amqp.ExchangeDirect,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return nil, fmt.Errorf("failed  amqp ExchangeDeclare err: %s", err)
	}

	return &MQ{channel: channel, Connection: connection}, nil
}

func (q *MQ) Publish(data []byte) error {
	err := q.channel.Publish(
		exchangeName,
		"",
		false,
		false,
		amqp.Publishing{
			Headers:      amqp.Table{},
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         data,
		})

	return err
}

