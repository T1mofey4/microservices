package v1

import (
	"context"
	"fmt"

	paymentV1 "github.com/T1mofey4/microservices/shared/pkg/proto/payment/v1"
	"github.com/google/uuid"

	"github.com/T1mofey4/microservices/order/internal/client/converter"
	"github.com/T1mofey4/microservices/order/internal/model"
)

func (c *paymentClient) PayOrder(ctx context.Context, userUUID, orderUUID uuid.UUID, method model.PaymentMethod) (uuid.UUID, error) {
	paymentMethod, err := converter.PaymentMethodToProto(method)
	if err != nil {
		return uuid.Nil, fmt.Errorf("convert payment method: %w", err)
	}

	req := &paymentV1.PayOrderRequest{
		UserUuid:      userUUID.String(),
		OrderUuid:     orderUUID.String(),
		PaymentMethod: paymentMethod,
	}

	res, err := c.client.PayOrder(ctx, req)
	if err != nil {
		return uuid.Nil, fmt.Errorf("payment request error: %w", err)
	}

	transactionUUID, err := uuid.Parse(res.TransactionUuid)
	if err != nil {
		return uuid.Nil, fmt.Errorf("convert uuid error: %w", err)
	}

	return transactionUUID, nil
}
