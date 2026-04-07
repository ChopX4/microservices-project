package config

type PaymentConfig interface {
	Address() string
}

type LoggerConfig interface {
	Level() string
	AsJson() bool
}
