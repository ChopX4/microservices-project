package main

import (
	"context"
	"fmt"
	"log"
	"syscall"
	"os/signal"

	"github.com/ChopX4/raketka/order/internal/app"
	"github.com/ChopX4/raketka/order/internal/config"
	"github.com/ChopX4/raketka/platform/pkg/closer"
)

const configPath = "./deploy/compose/order/.env"

func main() {
	if err := config.Load(configPath); err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

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
