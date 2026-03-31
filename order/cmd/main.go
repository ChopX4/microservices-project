package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
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
	"github.com/ChopX4/raketka/order/internal/config"
	migrator "github.com/ChopX4/raketka/order/internal/migrator"
	repo "github.com/ChopX4/raketka/order/internal/repository/order"
	orderService "github.com/ChopX4/raketka/order/internal/service/order"
	order_v1 "github.com/ChopX4/raketka/shared/pkg/openapi/order/v1"
	inventory_v1 "github.com/ChopX4/raketka/shared/pkg/proto/inventory/v1"
	payment_v1 "github.com/ChopX4/raketka/shared/pkg/proto/payment/v1"
)

const (
	shutdownTimeout = 10 * time.Second
	configPath      = "./deploy/compose/order/.env"
)

func main() {
	if err := config.Load(configPath); err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	paymentCon, err := grpc.NewClient(
		config.AppConfig().Payment.Address(),
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
		config.AppConfig().Inventory.Address(),
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

	db, err := sql.Open("pgx", config.AppConfig().Postgre.URI())
	if err != nil {
		log.Printf("failed to connect to database for migrations: %v\n", err)
		return
	}

	migratorRunner := migrator.NewMigrator(db, config.AppConfig().Postgre.MigrationsPath())
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

	conn, err := pgxpool.New(ctx, config.AppConfig().Postgre.URI())
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
		Addr:              config.AppConfig().Order.Address(),
		Handler:           r,
		ReadHeaderTimeout: config.AppConfig().Order.ReadTimeout(),
	}

	go func() {
		log.Printf("🚀 HTTP-сервер запущен на адресе %s\n", config.AppConfig().Order.Address())
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
