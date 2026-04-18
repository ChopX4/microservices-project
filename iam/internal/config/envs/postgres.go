package envs

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type postgreSQLEnvConfig struct {
	Host       string `env:"POSTGRES_HOST,required"`
	Port       string `env:"POSTGRES_PORT,required"`
	User       string `env:"POSTGRES_USER,required"`
	Password   string `env:"POSTGRES_PASSWORD,required"`
	Name       string `env:"POSTGRES_DB,required"`
	SSL        string `env:"POSTGRES_SSL_MODE,required"`
	Migrations string `env:"MIGRATION_DIRECTORY,required"`
}

type postgreSQLConfig struct {
	raw postgreSQLEnvConfig
}

func NewPostgreSQLConfig() (*postgreSQLConfig, error) {
	var raw postgreSQLEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &postgreSQLConfig{
		raw: raw,
	}, nil
}

func (c *postgreSQLConfig) URI() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.raw.User,
		c.raw.Password,
		c.raw.Host,
		c.raw.Port,
		c.raw.Name,
		c.raw.SSL,
	)
}

func (c *postgreSQLConfig) MigrationsPath() string {
	return c.raw.Migrations
}
