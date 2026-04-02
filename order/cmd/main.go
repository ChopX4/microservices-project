package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/ChopX4/raketka/order/internal/app"
	"github.com/ChopX4/raketka/order/internal/config"
	"github.com/ChopX4/raketka/platform/pkg/closer"
	"github.com/ChopX4/raketka/platform/pkg/logger"
)

const configPath = "./deploy/compose/order/.env"

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
