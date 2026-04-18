package envs

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type iamEnvConfig struct {
	Host string `env:"GRPC_HOST,required"`
	Port string `env:"GRPC_PORT,required"`
}

type iamConfig struct {
	raw iamEnvConfig
}

func NewIamConfig() (*iamConfig, error) {
	var raw iamEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &iamConfig{
		raw: raw,
	}, nil
}

func (c *iamConfig) Address() string {
	return net.JoinHostPort(c.raw.Host, c.raw.Port)
}
