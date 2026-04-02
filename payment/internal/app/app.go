package app

import (
	"context"
	"errors"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/ChopX4/raketka/payment/internal/config"
	"github.com/ChopX4/raketka/platform/pkg/closer"
	grpcHealth "github.com/ChopX4/raketka/platform/pkg/grpc/health"
	"github.com/ChopX4/raketka/platform/pkg/logger"
	payment_v1 "github.com/ChopX4/raketka/shared/pkg/proto/payment/v1"
)

type App struct {
	diContainer *diContainer
	grpcServer  *grpc.Server
	listener    net.Listener
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	return a.runGRPCServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.initListener,
		a.initGRPCServer,
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

func (a *App) initListener(_ context.Context) error {
	listener, err := net.Listen("tcp", config.AppConfig().Payment.Address())
	if err != nil {
		return err
	}

	closer.AddNamed("TCP listener", func(context.Context) error {
		if err := listener.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
			return err
		}

		return nil
	})

	a.listener = listener

	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer()
	closer.AddNamed("gRPC server", func(context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	reflection.Register(a.grpcServer)
	grpcHealth.RegisterService(a.grpcServer)
	payment_v1.RegisterPaymentServiceServer(a.grpcServer, a.diContainer.PaymentV1API(ctx))

	return nil
}

func (a *App) runGRPCServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🚀 gRPC PaymentService server listening on %s", config.AppConfig().Payment.Address()))

	return a.grpcServer.Serve(a.listener)
}
