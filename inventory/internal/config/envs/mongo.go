package envs

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type mongoEnvConfig struct {
	Host     string `env:"MONGO_HOST,required"`
	Port     string `env:"MONGO_PORT,required"`
	Name     string `env:"MONGO_DATABASE,required"`
	AuthDB   string `env:"MONGO_AUTH_DB,required"`
	Username string `env:"MONGO_INITDB_ROOT_USERNAME,required"`
	Password string `env:"MONGO_INITDB_ROOT_PASSWORD,required"`
}

type mongoConfig struct {
	raw mongoEnvConfig
}

func NewMongoConfig() (*mongoConfig, error) {
	var raw mongoEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &mongoConfig{
		raw: raw,
	}, nil
}

func (c *mongoConfig) URI() string {
	return fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/%s?authSource=%s",
		c.raw.Username,
		c.raw.Password,
		c.raw.Host,
		c.raw.Port,
		c.raw.Name,
		c.raw.AuthDB,
	)
}

func (c *mongoConfig) DbName() string {
	return c.raw.Name
}
