package envs

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type paymentEnvConfig struct {
	GrpcHost string `env:"GRPC_HOST,required"`
	GrpcPort string `env:"GRPC_PORT,required"`
}

type paymentConfig struct {
	raw paymentEnvConfig
}

func NewPaymentConfig() (*paymentConfig, error) {
	var raw paymentEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &paymentConfig{
		raw: raw,
	}, nil
}

func (c *paymentConfig) Address() string {
	return net.JoinHostPort(c.raw.GrpcHost, c.raw.GrpcPort)
}
