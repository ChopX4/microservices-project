package config

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type KafkaConfig interface {
	Brokers() []string
}

type AssembledConsumerConfig interface {
	Topic() string
	GroupID() string
}

type OrderConsumerConfig interface {
	Topic() string
	GroupID() string
}

type TelegramConfig interface {
	Token() string
}
