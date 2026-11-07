package payment

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/T1mofey4/microservices/payment/internal/model"
)

func (s *paymentService) PayOrder(ctx context.Context, orderUUID, userUUID string, method model.PaymentMethod) (string, error) {
	if method == model.PaymentMethodUnspecified {
		return "", model.ErrInvalidPaymentMethod
	}

	newUUID := uuid.New()
	log.Printf("Оплата прошла успешно, transaction_uuid: %s", newUUID.String())

	return newUUID.String(), nil
}
