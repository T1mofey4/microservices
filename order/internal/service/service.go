package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/T1mofey4/microservices/order/internal/model"
)

type OrderService interface {
	Create(ctx context.Context, userUUID uuid.UUID, filter model.PartsFilter) (model.Order, error)
	Get(ctx context.Context, orderUUID uuid.UUID) (model.Order, error)
	Pay(ctx context.Context, orderUUID uuid.UUID, method model.PaymentMethod) (uuid.UUID, error)
	Cancel(ctx context.Context, orderUUID uuid.UUID) error
}
