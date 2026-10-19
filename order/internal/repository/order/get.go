package order

import (
	"context"
	"fmt"

	"github.com/T1mofey4/microservices/order/internal/model"
	"github.com/T1mofey4/microservices/order/internal/repository/converter"
)

func (r *orderRepository) Get(ctx context.Context, orderUUID string) (model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	repoOrder, ok := r.store[orderUUID]
	if !ok {
		return model.Order{}, model.ErrOrderNotFound
	}

	modelOrder, err := converter.OrderToModel(repoOrder)
	if err != nil {
		return model.Order{}, fmt.Errorf("Get order %s: %w", orderUUID, err)
	}

	return modelOrder, nil
}
