package order

import (
	"sync"

	def "github.com/T1mofey4/microservices/order/internal/repository"
	repoModel "github.com/T1mofey4/microservices/order/internal/repository/model"
)

var _ def.OrderRepository = (*orderRepository)(nil)

type orderRepository struct {
	mu    sync.RWMutex
	store map[string]repoModel.Order
}

func NewOrderRepository() *orderRepository {
	return &orderRepository{store: make(map[string]repoModel.Order)}
}
