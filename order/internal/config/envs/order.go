package envs

import (
	"net"
	"time"

	"github.com/caarlos0/env/v11"
)

type orderEnvConfig struct {
	Host    string        `env:"HTTP_HOST,required"`
	Port    string        `env:"HTTP_PORT,required"`
	Timeout time.Duration `env:"HTTP_READ_TIMEOUT,required"`
}

type orderConfig struct {
	raw orderEnvConfig
}

func NewOrderConfig() (*orderConfig, error) {
	var raw orderEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderConfig{
		raw: raw,
	}, nil
}

func (c *orderConfig) Address() string {
	return net.JoinHostPort(c.raw.Host, c.raw.Port)
}

func (c *orderConfig) ReadTimeout() time.Duration {
	return c.raw.Timeout
}
