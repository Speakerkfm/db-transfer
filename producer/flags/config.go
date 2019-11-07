package flags

type Config struct {
	AmqpHost     string `long:"amqp-host" env:"AMQP_HOST" required:"true"`
	AmqpPort     int    `long:"amqp-port" env:"AMQP_PORT" required:"true"`
	AmqpUser     string `long:"amqp-user" env:"AMQP_USER" required:"true"`
	AmqpPassword string `long:"amqp-password" env:"AMQP_PASSWORD" required:"true"`

	RedisHost string `long:"redis-host" env:"REDIS_HOST" required:"true"`

	DatabaseHost     string `long:"database-host" env:"DATABASE_HOST" required:"true"`
	DatabasePort     string `long:"database-port" env:"DATABASE_PORT" default:"3306"`
	DatabaseName     string `long:"database-name" env:"DATABASE_NAME" required:"true"`
	DatabaseUser     string `long:"database-user" env:"DATABASE_USER" required:"true"`
	DatabasePassword string `long:"database-password" env:"DATABASE_PASSWORD" required:"true"`

	ExchangeName string `long:"exchange-name" env:"EXCHANGE_NAME" required:"true"`
	QueueName    string `long:"queue-name" env:"QUEUE_NAME" required:"true"`
	RoutingKey   string `long:"routing-key" env:"ROUTING_KEY" required:"false" default:""`
	DeadLetterExchangeName string `long:"dead-letter-exchange-name" env:"DEAD_LETTER_EXCHANGE_NAME" required:"true"`
	DeadLetterQueueName    string `long:"dead-letter-queue-name" env:"DEAD_LETTER_QUEUE_NAME" required:"true"`

	ServerHost     string `long:"server-host" env:"SERVER_HOST" required:"true"`

	SendType     string `long:"send-type" env:"SEND_TYPE" default:"tcp"`
}