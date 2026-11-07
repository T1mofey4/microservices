package service

import (
	"context"

	"github.com/T1mofey4/microservices/payment/internal/model"
)

type PaymentService interface {
	PayOrder(ctx context.Context, orderUUID, userUUID string, method model.PaymentMethod) (string, error)
}
