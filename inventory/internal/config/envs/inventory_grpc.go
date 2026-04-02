package envs

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type inventoryEnvGrpcConfig struct {
	GrpcHost string `env:"GRPC_HOST,required"`
	GprcPort string `env:"GRPC_PORT,required"`
}

type inventoryGrpcConfig struct {
	raw inventoryEnvGrpcConfig
}

func NewinventoryGrpcConfig() (*inventoryGrpcConfig, error) {
	var raw inventoryEnvGrpcConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &inventoryGrpcConfig{
		raw: raw,
	}, nil
}

func (c *inventoryGrpcConfig) Address() string {
	return net.JoinHostPort(c.raw.GrpcHost, c.raw.GprcPort)
}
