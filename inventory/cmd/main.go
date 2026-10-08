package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	inventoryV1 "github.com/T1mofey4/microservices/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const grpcPort = 50051

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	s := grpc.NewServer()

	storage := NewInventoryStorage(testParts())

	service := &inventoryService{storage: storage}

	inventoryV1.RegisterInventoryServiceServer(s, service)

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
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("Server stopped")
}

type inventoryService struct {
	inventoryV1.UnimplementedInventoryServiceServer
	storage *inventoryStorage
}

// GetPart отдаёт информацию о детали по её UUID
func (s *inventoryService) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	p, err := s.storage.GetPart(ctx, req.Uuid)
	if err != nil {
		switch {
		case errors.Is(err, ErrPartNotFound):
			return nil, status.Errorf(codes.NotFound, "деталь uuid: %s не найдена", req.Uuid)
		default:
			return nil, status.Error(codes.Internal, "внутренняя ошибка")
		}
	}

	return &inventoryV1.GetPartResponse{Part: p}, nil
}

// // ListParts отдаёт детали из storage, применяя фильтр из запроса. Пустой фильтр — все детали.
func (s *inventoryService) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	parts, err := s.storage.ListParts(ctx, req.Filter)
	if err != nil {
		return nil, status.Error(codes.Internal, "внутренняя ошибка")
	}

	return &inventoryV1.ListPartsResponse{Parts: parts}, nil
}
