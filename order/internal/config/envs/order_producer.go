package envs

import "github.com/caarlos0/env/v11"

type orderProducerEnvConfig struct {
	Topic string `env:"ORDER_PAID_TOPIC_NAME,required"`
}

type orderProducerConfig struct {
	raw orderProducerEnvConfig
}

func NewOrderProducerConfig() (*orderProducerConfig, error) {
	var raw orderProducerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderProducerConfig{
		raw: raw,
	}, nil
}

func (c *orderProducerConfig) Topic() string {
	return c.raw.Topic
}
