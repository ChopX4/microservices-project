package app

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderAPI "github.com/ChopX4/raketka/order/internal/api/order/v1"
	orderGRPCClients "github.com/ChopX4/raketka/order/internal/clients/grpc"
	orderInventoryClient "github.com/ChopX4/raketka/order/internal/clients/grpc/inventory"
	orderPaymentClient "github.com/ChopX4/raketka/order/internal/clients/grpc/payment"
	"github.com/ChopX4/raketka/order/internal/config"
	migrator "github.com/ChopX4/raketka/order/internal/migrator"
	"github.com/ChopX4/raketka/order/internal/repository"
	orderRepository "github.com/ChopX4/raketka/order/internal/repository/order"
	"github.com/ChopX4/raketka/order/internal/service"
	orderService "github.com/ChopX4/raketka/order/internal/service/order"
	"github.com/ChopX4/raketka/platform/pkg/closer"
	order_v1 "github.com/ChopX4/raketka/shared/pkg/openapi/order/v1"
	inventory_v1 "github.com/ChopX4/raketka/shared/pkg/proto/inventory/v1"
	payment_v1 "github.com/ChopX4/raketka/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	orderHandler order_v1.Handler

	orderService service.OrderService

	orderRepository repository.OrderRepository

	inventoryClient orderGRPCClients.InventoryClient
	paymentClient   orderGRPCClients.PaymentClient

	generatedInventoryClient inventory_v1.InventoryServiceClient
	generatedPaymentClient   payment_v1.PaymentServiceClient

	inventoryConn *grpc.ClientConn
	paymentConn   *grpc.ClientConn

	postgreSQLPool *pgxpool.Pool
}

func NewDIContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderHandler(ctx context.Context) order_v1.Handler {
	if d.orderHandler == nil {
		d.orderHandler = orderAPI.NewApi(d.OrderService(ctx))
	}

	return d.orderHandler
}

func (d *diContainer) OrderService(ctx context.Context) service.OrderService {
	if d.orderService == nil {
		d.orderService = orderService.NewService(
			d.OrderRepository(ctx),
			d.InventoryClient(ctx),
			d.PaymentClient(ctx),
		)
	}

	return d.orderService
}

func (d *diContainer) OrderRepository(ctx context.Context) repository.OrderRepository {
	if d.orderRepository == nil {
		d.orderRepository = orderRepository.NewRepository(d.PostgreSQLPool(ctx))
	}

	return d.orderRepository
}

func (d *diContainer) InventoryClient(ctx context.Context) orderGRPCClients.InventoryClient {
	if d.inventoryClient == nil {
		d.inventoryClient = orderInventoryClient.NewInventoryClient(d.GeneratedInventoryClient(ctx))
	}

	return d.inventoryClient
}

func (d *diContainer) PaymentClient(ctx context.Context) orderGRPCClients.PaymentClient {
	if d.paymentClient == nil {
		d.paymentClient = orderPaymentClient.NewPaymentClient(d.GeneratedPaymentClient(ctx))
	}

	return d.paymentClient
}

func (d *diContainer) GeneratedInventoryClient(ctx context.Context) inventory_v1.InventoryServiceClient {
	if d.generatedInventoryClient == nil {
		d.generatedInventoryClient = inventory_v1.NewInventoryServiceClient(d.InventoryConn(ctx))
	}

	return d.generatedInventoryClient
}

func (d *diContainer) GeneratedPaymentClient(ctx context.Context) payment_v1.PaymentServiceClient {
	if d.generatedPaymentClient == nil {
		d.generatedPaymentClient = payment_v1.NewPaymentServiceClient(d.PaymentConn(ctx))
	}

	return d.generatedPaymentClient
}

func (d *diContainer) InventoryConn(context.Context) *grpc.ClientConn {
	if d.inventoryConn == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().Inventory.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to connect to inventory grpc: %v", err))
		}

		closer.AddNamed("inventory gRPC connection", func(context.Context) error {
			return conn.Close()
		})

		d.inventoryConn = conn
	}

	return d.inventoryConn
}

func (d *diContainer) PaymentConn(context.Context) *grpc.ClientConn {
	if d.paymentConn == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().Payment.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to connect to payment grpc: %v", err))
		}

		closer.AddNamed("payment gRPC connection", func(context.Context) error {
			return conn.Close()
		})

		d.paymentConn = conn
	}

	return d.paymentConn
}

func (d *diContainer) RunMigrations() {
	db, err := sql.Open("pgx", config.AppConfig().Postgre.URI())
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database for migrations: %v", err))
	}

	migratorRunner := migrator.NewMigrator(db, config.AppConfig().Postgre.MigrationsPath())
	if err = migratorRunner.Up(); err != nil {
		_ = db.Close()
		panic(fmt.Sprintf("failed to migrate database: %v", err))
	}

	if err = db.Close(); err != nil {
		panic(fmt.Sprintf("failed to close migrations database connection: %v", err))
	}
}

func (d *diContainer) PostgreSQLPool(ctx context.Context) *pgxpool.Pool {
	if d.postgreSQLPool == nil {
		pool, err := pgxpool.New(ctx, config.AppConfig().Postgre.URI())
		if err != nil {
			panic(fmt.Sprintf("failed to create pgx pool: %v", err))
		}

		if err = pool.Ping(ctx); err != nil {
			pool.Close()
			panic(fmt.Sprintf("failed to ping postgres: %v", err))
		}

		closer.AddNamed("postgres pool", func(context.Context) error {
			pool.Close()
			return nil
		})

		d.postgreSQLPool = pool
	}

	return d.postgreSQLPool
}
