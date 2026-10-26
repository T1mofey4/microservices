package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/T1mofey4/microservices/order/internal/model"
)

func (s *orderService) Get(ctx context.Context, orderUUID uuid.UUID) (model.Order, error) {
	order, err := s.orderRepo.Get(ctx, orderUUID.String())
	if err != nil {
		return model.Order{}, fmt.Errorf("failed to get order %s: %w", orderUUID.String(), err)
	}

	return order, nil
}
