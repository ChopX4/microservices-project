package app

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	inventoryAPI "github.com/ChopX4/raketka/inventory/internal/api/inventory/v1"
	"github.com/ChopX4/raketka/inventory/internal/config"
	"github.com/ChopX4/raketka/inventory/internal/repository"
	inventoryRepository "github.com/ChopX4/raketka/inventory/internal/repository/part"
	"github.com/ChopX4/raketka/inventory/internal/service"
	inventoryService "github.com/ChopX4/raketka/inventory/internal/service/part"
	"github.com/ChopX4/raketka/platform/pkg/closer"
	"github.com/ChopX4/raketka/platform/pkg/logger"
	grpcAuth "github.com/ChopX4/raketka/platform/pkg/middleware/grpc"
	auth_v1 "github.com/ChopX4/raketka/shared/pkg/proto/auth/v1"
	inventory_v1 "github.com/ChopX4/raketka/shared/pkg/proto/inventory/v1"
)

type diContainer struct {
	inventoryV1API inventory_v1.InventoryServiceServer

	inventoryService service.InventoryService
	authInterceptor  *grpcAuth.AuthInterceptor

	inventoryRepository repository.InventoryRepository

	generatedIAMClient auth_v1.AuthServiceClient
	iamConn            *grpc.ClientConn

	mongoDBClient *mongo.Client
	mongoDBHandle *mongo.Database
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) InventoryV1API(ctx context.Context) inventory_v1.InventoryServiceServer {
	if d.inventoryV1API == nil {
		d.inventoryV1API = inventoryAPI.NewApi(d.InventoryService(ctx))
	}

	return d.inventoryV1API
}

func (d *diContainer) InventoryService(ctx context.Context) service.InventoryService {
	if d.inventoryService == nil {
		d.inventoryService = inventoryService.NewService(d.InventoryRepository(ctx))
	}

	return d.inventoryService
}

func (d *diContainer) AuthInterceptor(ctx context.Context) *grpcAuth.AuthInterceptor {
	if d.authInterceptor == nil {
		d.authInterceptor = grpcAuth.NewAuthInterceptor(d.GeneratedIAMClient(ctx))
	}

	return d.authInterceptor
}

func (d *diContainer) GeneratedIAMClient(ctx context.Context) auth_v1.AuthServiceClient {
	if d.generatedIAMClient == nil {
		d.generatedIAMClient = auth_v1.NewAuthServiceClient(d.IAMConn(ctx))
	}

	return d.generatedIAMClient
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

// InventoryRepository создает MongoDB-репозиторий с подготовленными индексами.
func (d *diContainer) InventoryRepository(ctx context.Context) repository.InventoryRepository {
	if d.inventoryRepository == nil {
		repo, err := inventoryRepository.NewRepository(ctx, d.MongoDBHandle(ctx))
		if err != nil {
			logger.Error(ctx, "failed to create inventory repository", zap.Error(err))
			panic(fmt.Sprintf("failed to create inventory repository: %v", err))
		}

		d.inventoryRepository = repo
	}

	return d.inventoryRepository
}

// MongoDBClient создает и проверяет клиент MongoDB.
func (d *diContainer) MongoDBClient(ctx context.Context) *mongo.Client {
	if d.mongoDBClient == nil {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
		if err != nil {
			logger.Error(ctx, "failed to connect to MongoDB", zap.Error(err))
			panic(fmt.Sprintf("failed to connect to MongoDB: %v", err))
		}

		if err = client.Ping(ctx, readpref.Primary()); err != nil {
			logger.Error(ctx, "failed to ping MongoDB", zap.Error(err))
			panic(fmt.Sprintf("failed to ping MongoDB: %v", err))
		}

		closer.AddNamed("MongoDB client", func(ctx context.Context) error {
			return client.Disconnect(ctx)
		})

		d.mongoDBClient = client
	}

	return d.mongoDBClient
}

func (d *diContainer) MongoDBHandle(ctx context.Context) *mongo.Database {
	if d.mongoDBHandle == nil {
		d.mongoDBHandle = d.MongoDBClient(ctx).Database(config.AppConfig().Mongo.DbName())
	}

	return d.mongoDBHandle
}
