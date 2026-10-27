package v1

import (
	"context"
	"errors"

	orderV1 "github.com/T1mofey4/microservices/shared/pkg/openapi/order/v1"

	"github.com/T1mofey4/microservices/order/internal/converter"
	"github.com/T1mofey4/microservices/order/internal/model"
)

func (a *api) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	filter := model.PartsFilter{UUIDs: req.PartUuids}
	order, err := a.orderService.Create(ctx, req.UserUUID, filter)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrPartNotFound):
			return &orderV1.NotFoundError{Code: 404, Message: "part not found"}, nil
		default:
			return &orderV1.InternalServerError{Code: 500, Message: "failed to create order"}, nil
		}
	}

	orderToDTO := converter.OrderToCreateOrderResponse(order)

	return &orderToDTO, nil
}
