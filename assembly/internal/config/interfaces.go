package config

import "time"

type KafkaConfig interface {
	Brokers() []string
}

type OrderPaidConsumerConfig interface {
	Topic() string
	GroupID() string
}

type OrderAssembledProducerConfig interface {
	Topic() string
}

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type MetricsConfig interface {
	CollectorEndpoint() string
	CollectorInterval() time.Duration
}
