package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/ChopX4/raketka/inventory/internal/config/envs"
)

var appConfig *config

type config struct {
	Iam       IamClientConfig
	Inventory InventoryConfig
	Logger    LoggerConfig
	Mongo     MongoConfig
}

func Load(path ...string) error {
	if err := godotenv.Load(path...); err != nil && !os.IsNotExist(err) {
		return err
	}

	iamCfg, err := envs.NewIAMClientConfig()
	if err != nil {
		return err
	}

	inventoryCfg, err := envs.NewinventoryGrpcConfig()
	if err != nil {
		return err
	}

	loggerCfg, err := envs.NewLoggerConfig()
	if err != nil {
		return err
	}

	mongoCfg, err := envs.NewMongoConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Iam:       iamCfg,
		Inventory: inventoryCfg,
		Logger:    loggerCfg,
		Mongo:     mongoCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
