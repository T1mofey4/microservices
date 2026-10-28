package v1

import (
	payV1 "github.com/T1mofey4/microservices/shared/pkg/proto/payment/v1"

	"github.com/T1mofey4/microservices/payment/internal/service"
)

var _ payV1.PaymentServiceServer = (*api)(nil)

type api struct {
	payV1.UnimplementedPaymentServiceServer
	service service.PaymentService
}

func NewApi(paymentService service.PaymentService) *api {
	return &api{
		service: paymentService,
	}
}
