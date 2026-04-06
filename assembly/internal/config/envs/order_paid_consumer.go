package envs

import "github.com/caarlos0/env/v11"

type orderPaidConsumerEnvConfig struct {
	Topic   string `env:"ORDER_PAID_TOPIC_NAME,required"`
	GroupID string `env:"ORDER_PAID_CONSUMER_GROUP_ID,required"`
}

type orderPaidConsumerConfig struct {
	raw orderPaidConsumerEnvConfig
}

func NewOrderPaidConsumerConfig() (*orderPaidConsumerConfig, error) {
	var raw orderPaidConsumerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderPaidConsumerConfig{
		raw: raw,
	}, nil
}

func (c *orderPaidConsumerConfig) Topic() string {
	return c.raw.Topic
}

func (c *orderPaidConsumerConfig) GroupID() string {
	return c.raw.GroupID
}
