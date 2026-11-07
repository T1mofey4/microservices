package order

import (
	"context"
	"fmt"

	"github.com/T1mofey4/microservices/order/internal/model"
	"github.com/T1mofey4/microservices/order/internal/repository"
	"github.com/T1mofey4/microservices/order/internal/repository/converter"
)

// Create создает заказ
func (r *orderRepository) Create(ctx context.Context, order model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := order.OrderUUID.String()

	if _, ok := r.store[id]; ok {
		return fmt.Errorf("create order %s: %w", id, repository.ErrOrderAlreadyExists)
	}

	repoOrder := converter.OrderToRepoModel(order)

	r.store[repoOrder.OrderUUID] = repoOrder

	return nil
}
