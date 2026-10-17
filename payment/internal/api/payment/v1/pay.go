package v1

import (
	"context"
	"errors"
	"log"

	payV1 "github.com/T1mofey4/microservices/shared/pkg/proto/payment/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/T1mofey4/microservices/payment/internal/converter"
	"github.com/T1mofey4/microservices/payment/internal/model"
)

func (a *api) PayOrder(ctx context.Context, req *payV1.PayOrderRequest) (*payV1.PayOrderResponse, error) {
	_, err := uuid.Parse(req.OrderUuid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid order_uuid")
	}

	_, err = uuid.Parse(req.UserUuid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_uuid")
	}

	method := converter.MethodToModel(req.PaymentMethod)

	transactionUUID, err := a.service.PayOrder(ctx, req.OrderUuid, req.UserUuid, method)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidPaymentMethod):
			return nil, status.Error(codes.InvalidArgument, "invalid payment_method")
		default:
			log.Printf("PayOrder failed: %v", err)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &payV1.PayOrderResponse{TransactionUuid: transactionUUID}, nil
}
