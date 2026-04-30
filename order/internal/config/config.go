package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/ChopX4/raketka/order/internal/config/envs"
)

var appConfig *config

type config struct {
	Iam               IamClientConfig
	Inventory         InventoryClientConfig
	Payment           PaymentClientConfig
	Order             OrderConfig
	Logger            LoggerConfig
	Postgre           PostgreSQLConfig
	Kafka             KafkaConfig
	AssembledConsumer AssembledConsumerConfig
	OrderProducer     OrderProducerConfig
	Metrics           MetricsConfig
}

func Load(paths ...string) error {
	if err := godotenv.Load(paths...); err != nil && !os.IsNotExist(err) {
		return err
	}

	iamCfg, err := envs.NewIamClientConfig()
	if err != nil {
		return err
	}

	inventoryCfg, err := envs.NewInventoryClientConfig()
	if err != nil {
		return err
	}

	paymentCfg, err := envs.NewPaymentClientConfig()
	if err != nil {
		return err
	}

	orderCfg, err := envs.NewOrderConfig()
	if err != nil {
		return err
	}

	loggerCfg, err := envs.NewLoggerConfig()
	if err != nil {
		return err
	}

	postgreCfg, err := envs.NewPostgreSQLConfig()
	if err != nil {
		return err
	}

	kafkaCfg, err := envs.NewKafkaConfig()
	if err != nil {
		return err
	}

	assembledConsumerCfg, err := envs.NewAssembledConsumerConfig()
	if err != nil {
		return err
	}

	orderProducerCfg, err := envs.NewOrderProducerConfig()
	if err != nil {
		return err
	}

	metricsCfg, err := envs.NewMetricsConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Iam:               iamCfg,
		Inventory:         inventoryCfg,
		Payment:           paymentCfg,
		Order:             orderCfg,
		Logger:            loggerCfg,
		Postgre:           postgreCfg,
		Kafka:             kafkaCfg,
		AssembledConsumer: assembledConsumerCfg,
		OrderProducer:     orderProducerCfg,
		Metrics:           metricsCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
