package config

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
