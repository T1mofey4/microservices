package order

import (
	"context"

	"github.com/T1mofey4/microservices/order/internal/model"
	"github.com/T1mofey4/microservices/order/internal/repository"
	"github.com/T1mofey4/microservices/order/internal/repository/converter"
)

// Update обновляет заказ
func (r *orderRepository) Update(ctx context.Context, orderUUID string, expectedStatus model.OrderStatus, updOrder model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	repoOrder, ok := r.store[orderUUID]
	if !ok {
		return model.ErrOrderNotFound
	}

	status := converter.StatusToRepoModel(expectedStatus)
	if repoOrder.Status != status {
		return repository.ErrStatusConflict
	}

	updRepoOrder := converter.OrderToRepoModel(updOrder)
	updRepoOrder.OrderUUID = orderUUID
	r.store[orderUUID] = updRepoOrder

	return nil
}
