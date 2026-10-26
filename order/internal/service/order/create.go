package order

import (
	"context"
	"fmt"
	"math"

	"github.com/google/uuid"

	"github.com/T1mofey4/microservices/order/internal/model"
)

// Create создает заказ, присваивает заказу uuid, считает общую сумму
func (s *orderService) Create(ctx context.Context, userUUID uuid.UUID, parts []uuid.UUID) (model.Order, error) {
	listParts, err := s.inventoryClient.ListParts(ctx, parts)
	if err != nil {
		return model.Order{}, fmt.Errorf("failed to get list of parts: %w", err)
	}

	var totalPrice float64
	for _, p := range listParts {
		totalPrice += p.Price
	}

	totalPrice = math.Round(totalPrice*100) / 100

	order := model.Order{
		OrderUUID:  uuid.New(),
		UserUUID:   userUUID,
		PartUUIDs:  parts,
		TotalPrice: totalPrice,
		Status:     model.StatusPendingPayment,
	}

	if err := s.orderRepo.Create(ctx, order); err != nil {
		return model.Order{}, fmt.Errorf("failed to create order: %w", err)
	}

	return order, nil
}
