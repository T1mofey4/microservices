package converter

import (
	"fmt"

	paymentV1 "github.com/T1mofey4/microservices/shared/pkg/proto/payment/v1"

	"github.com/T1mofey4/microservices/order/internal/model"
)

func PaymentMethodToProto(method model.PaymentMethod) (paymentV1.PaymentMethod, error) {
	switch method {
	case model.PaymentMethodUnknown:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED, fmt.Errorf("unknown payment method %s: %w", method, model.ErrInvalidPaymentMethod)
	case model.PaymentMethodCard:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CARD, nil
	case model.PaymentMethodSBP:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_SBP, nil
	case model.PaymentMethodCreditCard:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD, nil
	case model.PaymentMethodInvestorMoney:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY, nil
	default:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED, fmt.Errorf("unknown payment method %s: %w", method, model.ErrInvalidPaymentMethod)
	}
}
