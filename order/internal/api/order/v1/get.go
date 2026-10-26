package v1

import (
	"context"
	"errors"

	orderV1 "github.com/T1mofey4/microservices/shared/pkg/openapi/order/v1"

	"github.com/T1mofey4/microservices/order/internal/converter"
	"github.com/T1mofey4/microservices/order/internal/model"
)

// OrderByUUID получить заказ по UUID
func (s *api) OrderByUUID(ctx context.Context, orderUUID orderV1.OrderByUUIDParams) (orderV1.OrderByUUIDRes, error) {
	order, err := s.orderService.Get(ctx, orderUUID.OrderUUID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrOrderNotFound):
			return &orderV1.NotFoundError{Code: 404, Message: "order not found"}, nil
		default:
			return &orderV1.InternalServerError{Code: 500, Message: "failed to get order"}, nil
		}
	}

	orderToDTO := converter.OrderToOrderDto(order)

	return &orderToDTO, nil
}
