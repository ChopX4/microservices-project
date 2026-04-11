package config

type InventoryConfig interface {
	Address() string
}

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type MongoConfig interface {
	URI() string
	DbName() string
}
