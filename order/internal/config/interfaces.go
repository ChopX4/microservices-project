package config

import "time"

type InventoryClientConfig interface {
	Address() string
}

type IamClientConfig interface {
	Address() string
}

type PaymentClientConfig interface {
	Address() string
}

type OrderConfig interface {
	Address() string
	ReadTimeout() time.Duration
}

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type PostgreSQLConfig interface {
	URI() string
	MigrationsPath() string
}

type KafkaConfig interface {
	Brokers() []string
}

type AssembledConsumerConfig interface {
	Topic() string
	GroupID() string
}

type OrderProducerConfig interface {
	Topic() string
}
