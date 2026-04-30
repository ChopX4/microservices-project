package envs

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type metricsEnvConfig struct {
	CollectorEndpoint string        `env:"METRICS_COLLECTOR_ENDPOINT,required"`
	CollectorInterval time.Duration `env:"METRICS_COLLECTOR_INTERVAL,required"`
}

type metricsConfig struct {
	raw metricsEnvConfig
}

func NewMetricsConfig() (*metricsConfig, error) {
	var raw metricsEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &metricsConfig{
		raw: raw,
	}, nil
}

func (c *metricsConfig) CollectorEndpoint() string {
	return c.raw.CollectorEndpoint
}

func (c *metricsConfig) CollectorInterval() time.Duration {
	return c.raw.CollectorInterval
}
