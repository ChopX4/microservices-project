package envs

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type iamEnvClientConfig struct {
	Host string `env:"IAM_GRPC_HOST,required"`
	Port string `env:"IAM_GRPC_PORT,required"`
}

type iamClientConfig struct {
	raw iamEnvClientConfig
}

func NewIAMClientConfig() (*iamClientConfig, error) {
	var raw iamEnvClientConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &iamClientConfig{
		raw: raw,
	}, nil
}

func (c *iamClientConfig) Address() string {
	return net.JoinHostPort(c.raw.Host, c.raw.Port)
}
