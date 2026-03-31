package config

import "time"

type inventoryClientConfig interface {
	Address() string
}

type paymentClientConfig interface {
	Address() string
}

type orderConfig interface {
	Address() string
	ReadTimeout() time.Duration
}

type loggerConfig interface {
	Level() string
	AsJson() bool
}

type postgreSQLConfig interface {
	URI() string
	MigrationsPath() string
}
