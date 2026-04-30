package app

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderAPI "github.com/ChopX4/raketka/order/internal/api/order/v1"
	orderGRPCClients "github.com/ChopX4/raketka/order/internal/clients/grpc"
	orderInventoryClient "github.com/ChopX4/raketka/order/internal/clients/grpc/inventory"
	orderPaymentClient "github.com/ChopX4/raketka/order/internal/clients/grpc/payment"
	"github.com/ChopX4/raketka/order/internal/config"
	kafkaConverter "github.com/ChopX4/raketka/order/internal/converter/kafka"
	decoder "github.com/ChopX4/raketka/order/internal/converter/kafka/decoder"
	migrator "github.com/ChopX4/raketka/order/internal/migrator"
	"github.com/ChopX4/raketka/order/internal/repository"
	orderRepository "github.com/ChopX4/raketka/order/internal/repository/order"
	outboxRepository "github.com/ChopX4/raketka/order/internal/repository/outbox"
	"github.com/ChopX4/raketka/order/internal/service"
	assembledconsumer "github.com/ChopX4/raketka/order/internal/service/consumer/assembled_consumer"
	orderService "github.com/ChopX4/raketka/order/internal/service/order"
	orderMetrics "github.com/ChopX4/raketka/order/internal/service/order/metrics"
	outboxsender "github.com/ChopX4/raketka/order/internal/service/outbox"
	"github.com/ChopX4/raketka/platform/pkg/closer"
	kafkaConsumer "github.com/ChopX4/raketka/platform/pkg/kafka/consumer"
	"github.com/ChopX4/raketka/platform/pkg/logger"
	platformMetrics "github.com/ChopX4/raketka/platform/pkg/metrics"
	httpAuth "github.com/ChopX4/raketka/platform/pkg/middleware/http"
	"github.com/ChopX4/raketka/platform/pkg/pgxtx"
	order_v1 "github.com/ChopX4/raketka/shared/pkg/openapi/order/v1"
	auth_v1 "github.com/ChopX4/raketka/shared/pkg/proto/auth/v1"
	inventory_v1 "github.com/ChopX4/raketka/shared/pkg/proto/inventory/v1"
	payment_v1 "github.com/ChopX4/raketka/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	orderHandler   order_v1.Handler
	authMiddleware *httpAuth.AuthMiddleware

	orderService      service.OrderService
	assembledConsumer service.AssembledConsumer
	outboxSender      service.OutboxSender
	shipDecoder       kafkaConverter.ShipAssembledDecoder
	orderMetrics      *orderMetrics.Metrics

	orderRepository  repository.OrderRepository
	outboxRepository repository.OutboxRepository
	txManager        pgxtx.TxManager

	inventoryClient orderGRPCClients.InventoryClient
	paymentClient   orderGRPCClients.PaymentClient

	generatedIAMClient       auth_v1.AuthServiceClient
	generatedInventoryClient inventory_v1.InventoryServiceClient
	generatedPaymentClient   payment_v1.PaymentServiceClient

	iamConn       *grpc.ClientConn
	inventoryConn *grpc.ClientConn
	paymentConn   *grpc.ClientConn

	postgreSQLPool *pgxpool.Pool
	syncProducer   sarama.SyncProducer
	consumerGroup  sarama.ConsumerGroup
}

func NewDIContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) InitMetrics(ctx context.Context) error {
	if err := platformMetrics.InitProvider(
		ctx,
		config.AppConfig().Metrics.CollectorEndpoint(),
		config.AppConfig().Metrics.CollectorInterval(),
		"order",
	); err != nil {
		return err
	}

	closer.AddNamed("otel meter provider", platformMetrics.Shutdown)

	return nil
}

// OrderHandler лениво собирает HTTP/OpenAPI handler поверх сервисного слоя.
func (d *diContainer) OrderHandler(ctx context.Context) order_v1.Handler {
	if d.orderHandler == nil {
		d.orderHandler = orderAPI.NewApi(d.OrderService(ctx))
	}

	return d.orderHandler
}

func (d *diContainer) AuthMiddleware(ctx context.Context) *httpAuth.AuthMiddleware {
	if d.authMiddleware == nil {
		d.authMiddleware = httpAuth.NewAuthMiddleware(d.GeneratedIAMClient(ctx))
	}

	return d.authMiddleware
}

// OrderService связывает бизнес-логику с репозиторием и внешними gRPC-клиентами.
func (d *diContainer) OrderService(ctx context.Context) service.OrderService {
	if d.orderService == nil {
		d.orderService = orderService.NewService(
			d.OrderRepository(ctx),
			d.OutboxRepository(ctx),
			d.TxManager(ctx),
			d.InventoryClient(ctx),
			d.PaymentClient(ctx),
			config.AppConfig().OrderProducer.Topic(),
			d.OrderMetrics(ctx),
		)
	}

	return d.orderService
}

func (d *diContainer) AssembledConsumer(ctx context.Context) service.AssembledConsumer {
	if d.assembledConsumer == nil {
		d.assembledConsumer = assembledconsumer.NewAssembledConsumer(
			kafkaConsumer.NewConsumer(
				d.ConsumerGroup(ctx),
				[]string{config.AppConfig().AssembledConsumer.Topic()},
				logger.Logger(),
			),
			d.ShipDecoder(),
			d.OrderService(ctx),
		)
	}

	return d.assembledConsumer
}

func (d *diContainer) OutboxSender(ctx context.Context) service.OutboxSender {
	if d.outboxSender == nil {
		d.outboxSender = outboxsender.NewSender(
			d.OutboxRepository(ctx),
			d.SyncProducer(ctx),
			d.TxManager(ctx),
		)
	}

	return d.outboxSender
}

func (d *diContainer) ShipDecoder() kafkaConverter.ShipAssembledDecoder {
	if d.shipDecoder == nil {
		d.shipDecoder = decoder.NewShipDecoder()
	}

	return d.shipDecoder
}

func (d *diContainer) OrderMetrics(ctx context.Context) *orderMetrics.Metrics {
	if d.orderMetrics == nil {
		m, err := orderMetrics.New(ctx, otel.Meter("order.service"))
		if err != nil {
			logger.Error(ctx, "failed to initialize order metrics", zap.Error(err))
			panic(fmt.Sprintf("failed to initialize order metrics: %v", err))
		}

		d.orderMetrics = m
	}

	return d.orderMetrics
}

func (d *diContainer) OrderRepository(ctx context.Context) repository.OrderRepository {
	if d.orderRepository == nil {
		d.orderRepository = orderRepository.NewRepository(d.PostgreSQLPool(ctx))
	}

	return d.orderRepository
}

func (d *diContainer) OutboxRepository(ctx context.Context) repository.OutboxRepository {
	if d.outboxRepository == nil {
		d.outboxRepository = outboxRepository.NewOutboxRepository(d.PostgreSQLPool(ctx))
	}

	return d.outboxRepository
}

func (d *diContainer) TxManager(ctx context.Context) pgxtx.TxManager {
	if d.txManager == nil {
		d.txManager = pgxtx.NewTxManager(d.PostgreSQLPool(ctx))
	}

	return d.txManager
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

// GeneratedInventoryClient создает сгенерированный gRPC-клиент поверх общего соединения.
func (d *diContainer) GeneratedIAMClient(ctx context.Context) auth_v1.AuthServiceClient {
	if d.generatedIAMClient == nil {
		d.generatedIAMClient = auth_v1.NewAuthServiceClient(d.IAMConn(ctx))
	}

	return d.generatedIAMClient
}

// GeneratedInventoryClient создает сгенерированный gRPC-клиент поверх общего соединения.
func (d *diContainer) GeneratedInventoryClient(ctx context.Context) inventory_v1.InventoryServiceClient {
	if d.generatedInventoryClient == nil {
		d.generatedInventoryClient = inventory_v1.NewInventoryServiceClient(d.InventoryConn(ctx))
	}

	return d.generatedInventoryClient
}

// GeneratedPaymentClient создает сгенерированный gRPC-клиент поверх общего соединения.
func (d *diContainer) GeneratedPaymentClient(ctx context.Context) payment_v1.PaymentServiceClient {
	if d.generatedPaymentClient == nil {
		d.generatedPaymentClient = payment_v1.NewPaymentServiceClient(d.PaymentConn(ctx))
	}

	return d.generatedPaymentClient
}

func (d *diContainer) IAMConn(ctx context.Context) *grpc.ClientConn {
	if d.iamConn == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().Iam.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			logger.Error(ctx, "failed to create iam grpc connection", zap.Error(err))
			panic(fmt.Sprintf("failed to connect to iam grpc: %v", err))
		}

		closer.AddNamed("iam gRPC connection", func(context.Context) error {
			return conn.Close()
		})

		d.iamConn = conn
	}

	return d.iamConn
}

func (d *diContainer) InventoryConn(ctx context.Context) *grpc.ClientConn {
	if d.inventoryConn == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().Inventory.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			logger.Error(ctx, "failed to create inventory grpc connection", zap.Error(err))
			panic(fmt.Sprintf("failed to connect to inventory grpc: %v", err))
		}

		closer.AddNamed("inventory gRPC connection", func(context.Context) error {
			return conn.Close()
		})

		d.inventoryConn = conn
	}

	return d.inventoryConn
}

func (d *diContainer) PaymentConn(ctx context.Context) *grpc.ClientConn {
	if d.paymentConn == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().Payment.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			logger.Error(ctx, "failed to create payment grpc connection", zap.Error(err))
			panic(fmt.Sprintf("failed to connect to payment grpc: %v", err))
		}

		closer.AddNamed("payment gRPC connection", func(context.Context) error {
			return conn.Close()
		})

		d.paymentConn = conn
	}

	return d.paymentConn
}

// RunMigrations подключается к бд для миграций, после чего отключается.
// Дальше работа с идет через pxg Pool
func (d *diContainer) RunMigrations(ctx context.Context) {
	db, err := sql.Open("pgx", config.AppConfig().Postgre.URI())
	if err != nil {
		logger.Error(ctx, "failed to connect to database for migrations", zap.Error(err))
		panic(fmt.Sprintf("failed to connect to database for migrations: %v", err))
	}

	migratorRunner := migrator.NewMigrator(db, config.AppConfig().Postgre.MigrationsPath())
	if err = migratorRunner.Up(); err != nil {
		logger.Error(ctx, "failed to run database migrations", zap.Error(err))
		if closeErr := db.Close(); closeErr != nil {
			logger.Error(ctx, "failed to close migrations database connection after migration error", zap.Error(closeErr))
		}
		panic(fmt.Sprintf("failed to migrate database: %v", err))
	}

	if err = db.Close(); err != nil {
		logger.Error(ctx, "failed to close migrations database connection", zap.Error(err))
		panic(fmt.Sprintf("failed to close migrations database connection: %v", err))
	}
}

// PostgreSQLPool создает и проверяет рабочий pgx pool.
func (d *diContainer) PostgreSQLPool(ctx context.Context) *pgxpool.Pool {
	if d.postgreSQLPool == nil {
		pool, err := pgxpool.New(ctx, config.AppConfig().Postgre.URI())
		if err != nil {
			logger.Error(ctx, "failed to create postgres pool", zap.Error(err))
			panic(fmt.Sprintf("failed to create pgx pool: %v", err))
		}

		if err = pool.Ping(ctx); err != nil {
			logger.Error(ctx, "failed to ping postgres", zap.Error(err))
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

func (d *diContainer) SyncProducer(ctx context.Context) sarama.SyncProducer {
	if d.syncProducer == nil {
		cfg := sarama.NewConfig()
		cfg.Version = sarama.V3_6_0_0
		cfg.Producer.Return.Successes = true

		producer, err := sarama.NewSyncProducer(config.AppConfig().Kafka.Brokers(), cfg)
		if err != nil {
			logger.Error(ctx, "failed to create kafka sync producer", zap.Error(err))
			panic(fmt.Sprintf("failed to create kafka sync producer: %v", err))
		}

		closer.AddNamed("kafka sync producer", func(context.Context) error {
			return producer.Close()
		})

		d.syncProducer = producer
	}

	return d.syncProducer
}

func (d *diContainer) ConsumerGroup(ctx context.Context) sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		cfg := sarama.NewConfig()
		cfg.Version = sarama.V3_6_0_0
		cfg.Consumer.Return.Errors = true

		group, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().AssembledConsumer.GroupID(),
			cfg,
		)
		if err != nil {
			logger.Error(ctx, "failed to create kafka consumer group", zap.Error(err))
			panic(fmt.Sprintf("failed to create kafka consumer group: %v", err))
		}

		closer.AddNamed("kafka consumer group", func(context.Context) error {
			return group.Close()
		})

		d.consumerGroup = group
	}

	return d.consumerGroup
}
