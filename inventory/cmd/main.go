package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	inventoryApi "github.com/ChopX4/raketka/inventory/internal/api/inventory/v1"
	"github.com/ChopX4/raketka/inventory/internal/config"
	inventoryRepo "github.com/ChopX4/raketka/inventory/internal/repository/part"
	inventoryService "github.com/ChopX4/raketka/inventory/internal/service/part"
	inventory_v1 "github.com/ChopX4/raketka/shared/pkg/proto/inventory/v1"
)

const (
	configPath = "./deploy/compose/inventory/.env"
)

func main() {
	if err := config.Load(configPath); err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	lis, err := net.Listen("tcp", config.AppConfig().Inventory.Address())
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close: %v\n", cerr)
		}
	}()

	s := grpc.NewServer()

	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
	if err != nil {
		log.Printf("failed to connect to database: %v\n", err)
		return
	}
	defer func() {
		cerr := client.Disconnect(ctx)
		if cerr != nil {
			log.Printf("failed to disconnect: %v\n", cerr)
		}
	}()

	if err := client.Ping(ctx, nil); err != nil {
		log.Printf("failed to ping database: %v\n", err)
		return
	}

	db := client.Database(config.AppConfig().Mongo.DbName())

	repo, err := inventoryRepo.NewRepository(db)
	if err != nil {
		log.Printf("failed to create repo: %v\n", err)
		return
	}

	service := inventoryService.NewService(repo)
	api := inventoryApi.NewApi(service)

	inventory_v1.RegisterInventoryServiceServer(s, api)

	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC server listening on %s\n", config.AppConfig().Inventory.Address())
		err := s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
