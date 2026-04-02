package config

type paymentConfig interface {
	Address() string
}

type loggerConfig interface {
	Level() string
	AsJson() bool
}
