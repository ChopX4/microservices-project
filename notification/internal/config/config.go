package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/ChopX4/raketka/notification/internal/config/envs"
)

var appConfig *config

type config struct {
	Logger            LoggerConfig
	Kafka             KafkaConfig
	AssembledConsumer AssembledConsumerConfig
	OrderConsumer     OrderConsumerConfig
	TelegramConfig    TelegramConfig
}

func Load(paths ...string) error {
	if err := godotenv.Load(paths...); err != nil && !os.IsNotExist(err) {
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

	assembledCfg, err := envs.NewAssembledConsumerConfig()
	if err != nil {
		return err
	}

	orderCfg, err := envs.NewOrderConsumerConfig()
	if err != nil {
		return err
	}

	telegramCfg, err := envs.NewTelegramConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:            loggerCfg,
		Kafka:             kafkaCfg,
		AssembledConsumer: assembledCfg,
		OrderConsumer:     orderCfg,
		TelegramConfig:    telegramCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
