package converter

import (
	"log"

	paymentV1 "github.com/T1mofey4/microservices/shared/pkg/proto/payment/v1"

	"github.com/T1mofey4/microservices/payment/internal/model"
)

func MethodToModel(method paymentV1.PaymentMethod) model.PaymentMethod {
	switch method {
	case paymentV1.PaymentMethod_PAYMENT_METHOD_CARD:
		return model.PaymentMethodCard
	case paymentV1.PaymentMethod_PAYMENT_METHOD_SBP:
		return model.PaymentMethodSBP
	case paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return model.PaymentMethodCreditCard
	case paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return model.PaymentMethodInvestorMoney
	case paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED:
		return model.PaymentMethodUnspecified
	default:
		log.Printf("MethodToModel failed: invalid payment method: %d", method)
		return model.PaymentMethodUnspecified
	}
}
