package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	apiV1 "github.com/T1mofey4/microservices/inventory/internal/api/inventory/v1"
	repo "github.com/T1mofey4/microservices/inventory/internal/repository/part"
	service "github.com/T1mofey4/microservices/inventory/internal/service/part"
	inventoryV1 "github.com/T1mofey4/microservices/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50051

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

	dbURI := os.Getenv("MONGO_URI")

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	defer func() {
		dtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
		defer cancel()

		cerr := client.Disconnect(dtx)
		if cerr != nil {
			log.Printf("failed to disconnect: %v\n", cerr)
		}
	}()

	err = client.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	db := client.Database("inventory-service")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v", cerr)
		}
	}()

	s := grpc.NewServer()

	repo := repo.NewPartRepository(db)
	// напролняем базу тестовыми значениями
	if err = repo.InitParts(ctx); err != nil {
		return fmt.Errorf("init parts failed: %w", err)
	}

	service := service.NewService(repo)
	api := apiV1.NewApi(service)

	inventoryV1.RegisterInventoryServiceServer(s, api)

	// Включаем рефлексию для отладки
	reflection.Register(s)

	go func() {
		log.Printf("gRPC server listening on %d\n", grpcPort)
		err := s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// <-quit
	<-ctx.Done()
	log.Println("Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("Server stopped")

	return nil
}
