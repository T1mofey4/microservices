package v1

import (
	"context"
	"errors"

	orderV1 "github.com/T1mofey4/microservices/shared/pkg/openapi/order/v1"

	"github.com/T1mofey4/microservices/order/internal/converter"
	"github.com/T1mofey4/microservices/order/internal/model"
)

// OrderPay оплатить заказ
func (a *api) OrderPay(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.OrderPayParams) (orderV1.OrderPayRes, error) {
	paymentMethod, err := converter.PaymentMethodToModel(req.PaymentMethod)
	if err != nil {
		return &orderV1.BadRequestError{Code: 400, Message: "invalid payment method"}, nil
	}

	transactionUUID, err := a.orderService.Pay(ctx, params.OrderUUID, paymentMethod)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrOrderNotFound):
			return &orderV1.NotFoundError{Code: 404, Message: "order not found"}, nil
		case errors.Is(err, model.ErrOrderAlreadyPaid):
			return &orderV1.ConflictError{Code: 409, Message: "order already paid"}, nil
		case errors.Is(err, model.ErrOrderCancelled):
			return &orderV1.ConflictError{Code: 409, Message: "order cancelled"}, nil
		case errors.Is(err, model.ErrInvalidPaymentMethod):
			return &orderV1.BadRequestError{Code: 400, Message: "invalid payment method"}, nil
		case errors.Is(err, model.ErrInvalidArgument):
			return &orderV1.BadRequestError{Code: 400, Message: "invalid argument"}, nil
		default:
			return &orderV1.InternalServerError{Code: 500, Message: "internal server error"}, nil
		}
	}

	return &orderV1.PayOrderResponse{TransactionUUID: transactionUUID}, nil
}
