package envs

import "github.com/caarlos0/env/v11"

type orderConsumerEnvConfig struct {
	Topic   string `env:"ORDER_PAID_TOPIC_NAME,required"`
	GroupID string `env:"ORDER_PAID_CONSUMER_GROUP_ID,required"`
}

type orderConsumerConfig struct {
	raw orderConsumerEnvConfig
}

func NewOrderConsumerConfig() (*orderConsumerConfig, error) {
	var raw orderConsumerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderConsumerConfig{
		raw: raw,
	}, nil
}

func (c *orderConsumerConfig) Topic() string {
	return c.raw.Topic
}

func (c *orderConsumerConfig) GroupID() string {
	return c.raw.GroupID
}
