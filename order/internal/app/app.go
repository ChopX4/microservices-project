package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/ChopX4/raketka/order/internal/config"
	"github.com/ChopX4/raketka/platform/pkg/closer"
	"github.com/ChopX4/raketka/platform/pkg/logger"
	order_v1 "github.com/ChopX4/raketka/shared/pkg/openapi/order/v1"
)

const requestTimeout = 10 * time.Second

type App struct {
	diContainer *diContainer
	httpServer  *http.Server
}

// New создает приложение и последовательно поднимает его инфраструктурные зависимости.
func New(ctx context.Context) (*App, error) {
	a := &App{}

	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}

	return a, nil
}

// Run запускает HTTP-сервер и Kafka consumer приложения.
func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error, 3)

	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	closeResources := func() {
		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, requestTimeout)
		defer shutdownCancel()

		if closeErr := closer.CloseAll(shutdownCtx); closeErr != nil {
			logger.Error(ctx, "failed to close app resources", zap.Error(closeErr))
		}
	}

	go func() {
		if err := a.diContainer.AssembledConsumer(runCtx).RunAssembledConsumer(runCtx); err != nil {
			errCh <- fmt.Errorf("assembled consumer crashed: %w", err)
		}
	}()

	go func() {
		if err := a.diContainer.OutboxSender(runCtx).Run(runCtx); err != nil && !errors.Is(err, context.Canceled) {
			errCh <- fmt.Errorf("outbox sender crashed: %w", err)
		}
	}()

	go func() {
		if err := a.runHTTPServer(runCtx); err != nil {
			errCh <- fmt.Errorf("http server crashed: %w", err)
		}
	}()

	select {
	case <-runCtx.Done():
		closeResources()
		return runCtx.Err()
	case err := <-errCh:
		logger.Error(runCtx, "component crashed, shutting down", zap.Error(err))
		cancel()
		closeResources()
		return err
	}
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initMigrations,
		a.initHTTPServer,
	}

	for _, f := range inits {
		if err := f(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDIContainer()
	return nil
}

// initMigrations применяет миграции до создания рабочего runtime pool и старта HTTP-сервера.
func (a *App) initMigrations(ctx context.Context) error {
	a.diContainer.RunMigrations(ctx)
	return nil
}

// initHTTPServer собирает OpenAPI-хендлер, middleware и регистрирует graceful shutdown сервера.
func (a *App) initHTTPServer(ctx context.Context) error {
	orderServer, err := order_v1.NewServer(a.diContainer.OrderHandler(ctx))
	if err != nil {
		return err
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(requestTimeout))
	router.Use(a.diContainer.AuthMiddleware(ctx).Handle)
	router.Mount("/", orderServer)

	a.httpServer = &http.Server{
		Addr:        config.AppConfig().Order.Address(),
		Handler:     router,
		ReadTimeout: config.AppConfig().Order.ReadTimeout(),
	}

	closer.AddNamed("HTTP server", func(ctx context.Context) error {
		return a.httpServer.Shutdown(ctx)
	})

	return nil
}

func (a *App) runHTTPServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🚀 HTTP server listening on %s", config.AppConfig().Order.Address()))

	err := a.httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
