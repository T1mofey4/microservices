package repository

import (
	"context"
	"errors"

	"github.com/T1mofey4/microservices/order/internal/model"
)

var (
	ErrStatusConflict     = errors.New("order status conflict")
	ErrOrderAlreadyExists = errors.New("order already exists")
)

type OrderRepository interface {
	Create(ctx context.Context, order model.Order) error
	Get(ctx context.Context, orderUUID string) (model.Order, error)
	Update(ctx context.Context, orderUUID string, expectedStatus model.OrderStatus, updOrder model.Order) error
}
