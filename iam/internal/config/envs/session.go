package envs

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type sessionEnvConfig struct {
	TTL string `env:"SESSION_TTL,required"`
}

type sessionConfig struct {
	raw sessionEnvConfig
	ttl time.Duration
}

func NewSessionConfig() (*sessionConfig, error) {
	var raw sessionEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	ttl, err := time.ParseDuration(raw.TTL)
	if err != nil {
		return nil, err
	}

	return &sessionConfig{
		raw: raw,
		ttl: ttl,
	}, nil
}

func (c *sessionConfig) TTL() time.Duration {
	return c.ttl
}
