package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	orderOgenV1 "github.com/T1mofey4/microservices/shared/pkg/openapi/order/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1 "github.com/T1mofey4/microservices/order/internal/api/order/v1"
	inventoryV1 "github.com/T1mofey4/microservices/order/internal/client/grpc/inventory/v1"
	paymentV1 "github.com/T1mofey4/microservices/order/internal/client/grpc/payment/v1"
	orderRepo "github.com/T1mofey4/microservices/order/internal/repository/order"
	orderService "github.com/T1mofey4/microservices/order/internal/service/order"
)

const (
	httpPort          = "8080"
	inventoryPort     = "50051"
	paymentPort       = "50052"
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func main() {
	orderRepository := orderRepo.NewOrderRepository()

	inventoryAddr := net.JoinHostPort("localhost", inventoryPort)
	inventoryConn, err := grpc.NewClient(
		inventoryAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("ошибка создания клиента Inventory: %v", err)
	}
	inventoryClient := inventoryV1.NewInventoryClient(inventoryConn)
	defer func() {
		if err := inventoryConn.Close(); err != nil {
			log.Printf("Ошибка закрытия соединения Inventory: %v", err)
		}
	}()

	paymentAddr := net.JoinHostPort("localhost", paymentPort)
	paymentConn, err := grpc.NewClient(
		paymentAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("ошибка создания клиента Payment: %v", err)
	}
	paymentClient := paymentV1.NewPaymentClient(paymentConn)
	defer func() {
		if err := paymentConn.Close(); err != nil {
			log.Printf("Ошибка закрытия соединения Payment: %v", err)
		}
	}()

	service := orderService.NewOrderService(orderRepository, inventoryClient, paymentClient)
	orderAPI := orderV1.NewApi(service)
	orderServer, err := orderOgenV1.NewServer(orderAPI)
	if err != nil {
		log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
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

	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	// Graceful shutdown
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch
	log.Println("🛑 Завершение работы сервера...")
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}
	log.Println("✅ Сервер остановлен")
}
