package config

type PaymentConfig interface {
	Address() string
}

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type KafkaConfig interface {
	Brokers() []string
}

type assembledConsumerConfig interface {
	Topic() string
	GroupID() string
}

type OrderProducerConfig interface {
	Topic() string
}
