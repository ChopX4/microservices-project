package envs

import "github.com/caarlos0/env/v11"

type orderAssembledProducerEnvConfig struct {
	Topic string `env:"ORDER_ASSEMBLED_TOPIC_NAME,required"`
}

type orderAssembledProducerConfig struct {
	raw orderAssembledProducerEnvConfig
}

func NewOrderAssembledProducerConfig() (*orderAssembledProducerConfig, error) {
	var raw orderAssembledProducerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderAssembledProducerConfig{
		raw: raw,
	}, nil
}

func (c *orderAssembledProducerConfig) Topic() string {
	return c.raw.Topic
}
