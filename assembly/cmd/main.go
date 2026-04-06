package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/ChopX4/raketka/assembly/internal/config"
	"github.com/ChopX4/raketka/assembly/internal/app"
	"github.com/ChopX4/raketka/platform/pkg/closer"
	"github.com/ChopX4/raketka/platform/pkg/logger"
	"go.uber.org/zap"
)

var (
	configPath = "./deploy/compose/assembly/.env"
)

func main() {
	if err := config.Load(configPath); err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	if err := logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
	); err != nil {
		panic(fmt.Errorf("failed to init logger: %w", err))
	}

	closer.SetLogger(logger.Logger())

	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	application, err := app.New(appCtx)
	if err != nil {
		logger.Error(appCtx, "failed to init app", zap.Error(err))
		return
	}

	if err = application.Run(appCtx); err != nil {
		logger.Error(appCtx, "failed to run app", zap.Error(err))
	}
}
