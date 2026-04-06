package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/ChopX4/raketka/payment/internal/app"
	"github.com/ChopX4/raketka/payment/internal/config"
	"github.com/ChopX4/raketka/platform/pkg/closer"
	"github.com/ChopX4/raketka/platform/pkg/logger"
)

const configPath = "./deploy/compose/payment/.env"

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
		log.Printf("failed to init app: %v\n", err)
		return
	}

	if err = application.Run(appCtx); err != nil {
		log.Printf("failed to run app: %v\n", err)
	}
}
