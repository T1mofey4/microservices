package main

import (
	"log"
	"sync"
	"time"

	orderV1 "github.com/T1mofey4/microservices/week_1/shared/pkg/openapi/order/v1"
)

const (
	httpPort          = "8080"
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func main() {
	storage := NewOrderStorage()

	orderHandler := NewOrderHandler(storage)

	_, err := orderV1.NewServer(orderHandler)
	if err != nil {
		log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
	}

}

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*orderV1.OrderDto
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*orderV1.OrderDto),
	}
}

type OrderHandler struct {
	storage *OrderStorage
}

func NewOrderHandler(storage *OrderStorage) *OrderHandler {
	return &OrderHandler{
		storage: storage,
	}
}
