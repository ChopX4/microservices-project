package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/ChopX4/raketka/iam/internal/config/envs"
)

var appConfig *config

type config struct {
	Iam     IamConfig
	Logger  LoggerConfig
	Postgre PostgreSQLConfig
	Redis   RedisConfig
	Session SessionConfig
}

func Load(paths ...string) error {
	if err := godotenv.Load(paths...); err != nil && !os.IsNotExist(err) {
		return err
	}

	iamCfg, err := envs.NewIamConfig()
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

	redisCfg, err := envs.NewRedisConfig()
	if err != nil {
		return err
	}

	sessionCfg, err := envs.NewSessionConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Iam:     iamCfg,
		Logger:  loggerCfg,
		Postgre: postgreCfg,
		Redis:   redisCfg,
		Session: sessionCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
