package payment

import def "github.com/T1mofey4/microservices/payment/internal/service"

var _ def.PaymentService = (*paymentService)(nil)

type paymentService struct{}

func NewPaymentService() *paymentService {
	return &paymentService{}
}
