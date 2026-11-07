package grpc

import (
	"context"

	"github.com/google/uuid"

	"github.com/T1mofey4/microservices/order/internal/model"
)

type InventoryClient interface {
	ListParts(ctx context.Context, filter model.PartsFilter) ([]model.PartInfo, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, userUUID, orderUUID uuid.UUID, method model.PaymentMethod) (uuid.UUID, error)
}
