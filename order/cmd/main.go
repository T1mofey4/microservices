package main

import (
	"context"
	"errors"
	"fmt"
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
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1 "github.com/T1mofey4/microservices/order/internal/api/order/v1"
	inventoryV1 "github.com/T1mofey4/microservices/order/internal/client/grpc/inventory/v1"
	paymentV1 "github.com/T1mofey4/microservices/order/internal/client/grpc/payment/v1"
	"github.com/T1mofey4/microservices/order/internal/migrator"
	orderRepo "github.com/T1mofey4/microservices/order/internal/repository/order"
	orderService "github.com/T1mofey4/microservices/order/internal/service/order"
	orderOgenV1 "github.com/T1mofey4/microservices/shared/pkg/openapi/order/v1"
)

const (
	httpPort          = "8080"
	inventoryPort     = "50051"
	paymentPort       = "50052"
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	err := godotenv.Load("../.env")
	if err != nil {
		return fmt.Errorf("failed to load .env file: %w", err)
	}

	dbURI := os.Getenv("DB_URI")
	if dbURI == "" {
		return fmt.Errorf("DB_URI is not set")
	}

	pool, err := pgxpool.New(ctx, dbURI)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("database is unavailable: %w", err)
	}

	db := stdlib.OpenDBFromPool(pool)
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("failed to close db: %v", err)
		}
	}()

	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	if migrationsDir == "" {
		migrationsDir = "migrations"
	}
	migratorRunner := migrator.NewMigrator(db, migrationsDir)

	err = migratorRunner.Up()
	if err != nil {
		return fmt.Errorf("database migration error: %w", err)
	}

	orderRepository := orderRepo.NewOrderRepository(pool)

	inventoryAddr := net.JoinHostPort("localhost", inventoryPort)
	inventoryConn, err := grpc.NewClient(
		inventoryAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("ошибка создания клиента Inventory: %w", err)
	}
	defer func() {
		if err := inventoryConn.Close(); err != nil {
			log.Printf("Ошибка закрытия соединения Inventory: %v", err)
		}
	}()
	inventoryClient := inventoryV1.NewInventoryClient(inventoryConn)

	paymentAddr := net.JoinHostPort("localhost", paymentPort)
	paymentConn, err := grpc.NewClient(
		paymentAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("ошибка создания клиента Payment: %w", err)
	}
	defer func() {
		if err := paymentConn.Close(); err != nil {
			log.Printf("Ошибка закрытия соединения Payment: %v", err)
		}
	}()
	paymentClient := paymentV1.NewPaymentClient(paymentConn)

	service := orderService.NewOrderService(orderRepository, inventoryClient, paymentClient)
	orderAPI := orderV1.NewApi(service)
	orderServer, err := orderOgenV1.NewServer(orderAPI)
	if err != nil {
		return fmt.Errorf("ошибка создания сервера OpenAPI: %w", err)
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

	defer func() {
		if err := server.Close(); err != nil {
			log.Printf("Ошибка закрытия сервера: %v", err)
		}
	}()

	chErr := make(chan error, 1)
	ln, err := net.Listen("tcp", server.Addr)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s", httpPort)
		err := server.Serve(ln)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			chErr <- err
			log.Printf("❌ Ошибка запуска сервера: %v", err)
		}
	}()

	// Graceful shutdown
	select {
	case <-ctx.Done():
		log.Println("🛑 Завершение работы сервера...")

		ctxShutdown, cancel := context.WithTimeout(context.WithoutCancel(ctx), shutdownTimeout)
		defer cancel()

		err = server.Shutdown(ctxShutdown)
		if err != nil {
			log.Printf("❌ Ошибка при остановке сервера: %v", err)
		}
		log.Println("✅ Сервер остановлен")
	case srvErr := <-chErr:
		return fmt.Errorf("http server failed: %w", srvErr)
	}

	return nil
}
