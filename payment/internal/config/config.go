package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/ChopX4/raketka/payment/internal/config/envs"
)

var appConfig *config

type config struct {
	Payment PaymentConfig
	Logger  LoggerConfig
}

func Load(paths ...string) error {
	if err := godotenv.Load(paths...); err != nil && !os.IsNotExist(err) {
		return err
	}

	paymentCfg, err := envs.NewPaymentConfig()
	if err != nil {
		return err
	}

	loggerCfg, err := envs.NewLoggerConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Payment: paymentCfg,
		Logger:  loggerCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
