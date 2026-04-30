package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/ChopX4/raketka/assembly/internal/config/envs"
)

var appConfig *config

type config struct {
	Kafka                  KafkaConfig
	OrderPaidConsumer      OrderPaidConsumerConfig
	OrderAssembledProducer OrderAssembledProducerConfig
	Logger                 LoggerConfig
	Metrics                MetricsConfig
}

func Load(path ...string) error {
	if err := godotenv.Load(path...); err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerCfg, err := envs.NewLoggerConfig()
	if err != nil {
		return err
	}

	kafkaCfg, err := envs.NewKafkaConfig()
	if err != nil {
		return err
	}

	orderPaidConsumerCfg, err := envs.NewOrderPaidConsumerConfig()
	if err != nil {
		return err
	}

	orderAssembledProducerCfg, err := envs.NewOrderAssembledProducerConfig()
	if err != nil {
		return err
	}

	metricsCfg, err := envs.NewMetricsConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Kafka:                  kafkaCfg,
		OrderPaidConsumer:      orderPaidConsumerCfg,
		OrderAssembledProducer: orderAssembledProducerCfg,
		Logger:                 loggerCfg,
		Metrics:                metricsCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
