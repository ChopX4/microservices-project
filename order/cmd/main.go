package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderApi "github.com/ChopX4/raketka/order/internal/api/order/v1"
	orderInventoryClient "github.com/ChopX4/raketka/order/internal/clients/grpc/inventory"
	orderPaymentClient "github.com/ChopX4/raketka/order/internal/clients/grpc/payment"
	migrator "github.com/ChopX4/raketka/order/internal/migrator"
	repo "github.com/ChopX4/raketka/order/internal/repository/order"
	orderService "github.com/ChopX4/raketka/order/internal/service/order"
	order_v1 "github.com/ChopX4/raketka/shared/pkg/openapi/order/v1"
	inventory_v1 "github.com/ChopX4/raketka/shared/pkg/proto/inventory/v1"
	payment_v1 "github.com/ChopX4/raketka/shared/pkg/proto/payment/v1"
)

const (
	httpPort           = "8080"
	paymentAddress     = "localhost:50052"
	inventoryAddress   = "localhost:50051"
	dbURI              = "postgres://order-service-user:order-service-password@localhost:5432/order-service?sslmode=disable"
	orderMigrationsDir = "./order/migrations"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func main() {
	paymentCon, err := grpc.NewClient(
		paymentAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return
	}
	defer func() {
		if cerr := paymentCon.Close(); cerr != nil {
			log.Printf("failed to close connect: %v", cerr)
		}
	}()

	inventoryCon, err := grpc.NewClient(
		inventoryAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return
	}
	defer func() {
		if cerr := inventoryCon.Close(); cerr != nil {
			log.Printf("failed to close connect: %v", cerr)
		}
	}()

	paymentClient := payment_v1.NewPaymentServiceClient(paymentCon)
	inventoryClient := inventory_v1.NewInventoryServiceClient(inventoryCon)
	payServiceClient := orderPaymentClient.NewPaymentClient(paymentClient)
	invServiceClient := orderInventoryClient.NewInventoryClient(inventoryClient)

	ctx := context.Background()

	db, err := sql.Open("pgx", dbURI)
	if err != nil {
		log.Printf("failed to connect to database for migrations: %v\n", err)
		return
	}

	migratorRunner := migrator.NewMigrator(db, orderMigrationsDir)
	err = migratorRunner.Up()
	if err != nil {
		if cerr := db.Close(); cerr != nil {
			log.Printf("Ошибка закрытия подключения к базе данных: %v\n", err)
		}

		log.Printf("Ошибка миграции базы данных: %v\n", err)
		return
	}
	if cerr := db.Close(); cerr != nil {
		log.Printf("Ошибка закрытия подключения к базе данных: %v\n", err)
	}

	conn, err := pgxpool.New(ctx, dbURI)
	if err != nil {
		log.Printf("failed to connect to database: %v\n", err)
		return
	}
	defer conn.Close()

	if err := conn.Ping(ctx); err != nil {
		log.Printf("База данных недоступна: %v\n", err)
		return
	}

	repo := repo.NewRepository(conn)
	s := orderService.NewService(repo, invServiceClient, payServiceClient)
	api := orderApi.NewApi(s)

	orderServer, err := order_v1.NewServer(api)
	if err != nil {
		log.Printf("ошибка создания сервера OpenAPI: %v", err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err := server.ListenAndServe()
		if err != nil {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
