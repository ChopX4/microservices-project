package config

import (
	"os"

	"github.com/ChopX4/raketka/order/internal/config/envs"
	"github.com/joho/godotenv"
)

var appConfig *config

type config struct {
	Inventory inventoryClientConfig
	Payment   paymentClientConfig
	Order     orderConfig
	Logger    loggerConfig
	Postgre   postgreSQLConfig
}

func Load(paths ...string) error {
	if err := godotenv.Load(paths...); err != nil && !os.IsNotExist(err) {
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

	appConfig = &config{
		Inventory: inventoryCfg,
		Payment:   paymentCfg,
		Order:     orderCfg,
		Logger:    loggerCfg,
		Postgre:   postgreCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
