package services

import (
	flagsCfg "db-transfer/consumer/pkg/flags"
	"db-transfer/consumer/pkg/log"
	"fmt"
	"github.com/streadway/amqp"
	"time"
)

type MqConsumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	consume func(amqp.Delivery)
	config  *flagsCfg.Config
}

func NewMqConsumer(config *flagsCfg.Config, consume func(amqp.Delivery)) *MqConsumer {
	c := &MqConsumer{
		consume: consume,
		config:  config,
	}

	amqpURI := fmt.Sprintf("amqp://%s:%s@%s:%d/", config.AmqpUser, config.AmqpPassword, config.AmqpHost, config.AmqpPort)

	c.conn = connectToRabbitMQ(amqpURI)

	log.Info().Msg("RabbitMQ: getting Channel")
	var err error
	c.channel, err = c.conn.Channel()
	if err != nil {
		panic(err)
	}

	exchangeDeclare(c.channel, config)
	announceQueue(c.channel, config)

	deliveries, err := c.channel.Consume(
		config.QueueName,
		c.tag,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	notify := c.conn.NotifyClose(make(chan *amqp.Error))
	go mqHandle(deliveries, c.consume, notify)

	return c
}

func (c *MqConsumer) Shutdown() error {
	// will close() the deliveries channel
	if err := c.channel.Cancel(c.tag, true); err != nil {
		return fmt.Errorf("RabbitMQ: MqConsumer cancel failed, err %s", err)
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("RabbitMQ: connection close error %s", err)
	}

	log.Info().Msg("RabbitMQ: AMQP shutdown OK")

	return nil
}

func mqHandle(deliveries <-chan amqp.Delivery, consume func(amqp.Delivery), errorChan <-chan *amqp.Error) {
	log.Info().Msg("RabbitMQ: handler starting")
	for {
		select {
		case <-errorChan:
			log.Info().Msg("RabbitMQ: handler stop")
			return
		case d, ok := <-deliveries:
			if ok {
				consume(d)
			}
		}
	}
}

func connectToRabbitMQ(uri string) *amqp.Connection {
	for {
		log.Info().Msg("RabbitMQ: trying create connect")
		conn, err := amqp.Dial(uri)

		if err == nil {
			return conn
		}

		log.Warn().Err(err).Msg("RabbitMQ: connect creation error")
		time.Sleep(500 * time.Millisecond)
	}
}

func exchangeDeclare(channel *amqp.Channel, config *flagsCfg.Config) {
	log.Info().Msg("RabbitMQ: exchange declare")
	if err := channel.ExchangeDeclare(
		config.DeadLetterExchangeName,
		amqp.ExchangeFanout,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		panic(err)
	}

	if err := channel.ExchangeDeclare(
		config.ExchangeName,
		amqp.ExchangeDirect,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		panic(err)
	}
}

func announceQueue(channel *amqp.Channel, config *flagsCfg.Config) {
	log.Info().Msg("RabbitMQ: announce queue")

	//queues declare
	dlQueue, err := channel.QueueDeclare(
		config.DeadLetterExchangeName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	queue, err := channel.QueueDeclare(
		config.QueueName,
		true,
		false,
		false,
		false,
		amqp.Table{"x-dead-letter-exchange": config.DeadLetterExchangeName},
	)
	if err != nil {
		panic(err)
	}

	//queues binding
	log.Info().Msg("RabbitMQ: queues binding")

	if err := channel.QueueBind(
		dlQueue.Name,
		config.RoutingKey,
		config.DeadLetterExchangeName,
		false,
		nil,
	); err != nil {
		panic(err)
	}

	if err := channel.QueueBind(
		queue.Name,
		config.RoutingKey,
		config.ExchangeName,
		false,
		nil,
	); err != nil {
		panic(err)
	}
}
