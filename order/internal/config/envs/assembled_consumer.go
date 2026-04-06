package envs

import "github.com/caarlos0/env/v11"

type assembledConsumerEnvConfig struct {
	Topic   string `env:"ORDER_ASSEMBLED_TOPIC_NAME,required"`
	GroupID string `env:"ORDER_ASSEMBLED_CONSUMER_GROUP_ID,required"`
}

type assembledConsumerConfig struct {
	raw assembledConsumerEnvConfig
}

func NewAssembledConsumerConfig() (*assembledConsumerConfig, error) {
	var raw assembledConsumerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &assembledConsumerConfig{
		raw: raw,
	}, nil
}

func (c *assembledConsumerConfig) Topic() string {
	return c.raw.Topic
}

func (c *assembledConsumerConfig) GroupID() string {
	return c.raw.GroupID
}
