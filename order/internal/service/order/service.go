package order

import (
	clients "github.com/T1mofey4/microservices/order/internal/client/grpc"
	repo "github.com/T1mofey4/microservices/order/internal/repository"
	def "github.com/T1mofey4/microservices/order/internal/service"
)

var _ def.OrderService = (*orderService)(nil)

type orderService struct {
	orderRepo       repo.OrderRepository
	inventoryClient clients.InventoryClient
	paymentClient   clients.PaymentClient
}

func NewOrderService(repository repo.OrderRepository, inventoryClient clients.InventoryClient, paymentClient clients.PaymentClient) *orderService {
	return &orderService{
		orderRepo:       repository,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}
