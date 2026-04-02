package envs

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type inventoryEnvClientConfig struct {
	Host string `env:"INVENTORY_GRPC_HOST,required"`
	Port string `env:"INVENTORY_GRPC_PORT,required"`
}

type inventoryClientConfig struct {
	raw inventoryEnvClientConfig
}

func NewInventoryClientConfig() (*inventoryClientConfig, error) {
	var raw inventoryEnvClientConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &inventoryClientConfig{
		raw: raw,
	}, nil
}

func (c *inventoryClientConfig) Address() string {
	return net.JoinHostPort(c.raw.Host, c.raw.Port)
}
