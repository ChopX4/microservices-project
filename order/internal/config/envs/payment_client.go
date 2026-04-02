package envs

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type paymentEnvClientConfig struct {
	Host string `env:"PAYMENT_GRPC_HOST,required"`
	Port string `env:"PAYMENT_GRPC_PORT,required"`
}

type paymentClientConfig struct {
	raw paymentEnvClientConfig
}

func NewPaymentClientConfig() (*paymentClientConfig, error) {
	var raw paymentEnvClientConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &paymentClientConfig{
		raw: raw,
	}, nil
}

func (c *paymentClientConfig) Address() string {
	return net.JoinHostPort(c.raw.Host, c.raw.Port)
}
