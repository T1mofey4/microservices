package v1

import (
	"context"
	"errors"

	orderV1 "github.com/T1mofey4/microservices/shared/pkg/openapi/order/v1"

	"github.com/T1mofey4/microservices/order/internal/model"
)

func (a *api) OrderCancel(ctx context.Context, params orderV1.OrderCancelParams) (orderV1.OrderCancelRes, error) {
	if err := a.orderService.Cancel(ctx, params.OrderUUID); err != nil {
		switch {
		case errors.Is(err, model.ErrOrderNotFound):
			return &orderV1.NotFoundError{Code: 404, Message: "order not found"}, nil
		case errors.Is(err, model.ErrOrderAlreadyPaid):
			return &orderV1.ConflictError{Code: 409, Message: "order already paid"}, nil
		case errors.Is(err, model.ErrOrderCancelled):
			return &orderV1.ConflictError{Code: 409, Message: "order already cancelled"}, nil
		case errors.Is(err, model.ErrInvalidArgument):
			return &orderV1.BadRequestError{Code: 400, Message: "invalid argument"}, nil
		default:
			return &orderV1.InternalServerError{Code: 500, Message: "internal error"}, nil
		}
	}

	return &orderV1.OrderCancelNoContent{}, nil
}
