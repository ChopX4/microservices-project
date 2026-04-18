package envs

import (
	"net"
	"time"

	"github.com/caarlos0/env/v11"
)

type redisEnvConfig struct {
	Host        string `env:"REDIS_HOST,required"`
	Port        string `env:"REDIS_PORT,required"`
	Timeout     string `env:"REDIS_CONNECTION_TIMEOUT,required"`
	MaxIdle     int    `env:"REDIS_MAX_IDLE,required"`
	IdleTimeout string `env:"REDIS_IDLE_TIMEOUT,required"`
}

type redisConfig struct {
	raw               redisEnvConfig
	connectionTimeout time.Duration
	idleTimeout       time.Duration
}

func NewRedisConfig() (*redisConfig, error) {
	var raw redisEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	connectionTimeout, err := time.ParseDuration(raw.Timeout)
	if err != nil {
		return nil, err
	}

	idleTimeout, err := time.ParseDuration(raw.IdleTimeout)
	if err != nil {
		return nil, err
	}

	return &redisConfig{
		raw:               raw,
		connectionTimeout: connectionTimeout,
		idleTimeout:       idleTimeout,
	}, nil
}

func (c *redisConfig) Address() string {
	return net.JoinHostPort(c.raw.Host, c.raw.Port)
}

func (c *redisConfig) ConnectionTimeout() time.Duration {
	return c.connectionTimeout
}

func (c *redisConfig) MaxIdle() int {
	return c.raw.MaxIdle
}

func (c *redisConfig) IdleTimeout() time.Duration {
	return c.idleTimeout
}
