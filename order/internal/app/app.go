package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

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

// Run запускает HTTP-сервер приложения.
func (a *App) Run(ctx context.Context) error {
	return a.runHTTPServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
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

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
	)
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
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
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
