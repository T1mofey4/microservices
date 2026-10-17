package v1

import (
	payV1 "github.com/T1mofey4/microservices/shared/pkg/proto/payment/v1"

	"github.com/T1mofey4/microservices/payment/internal/service"
)

var _ payV1.PaymentServiceServer = (*api)(nil)

type api struct {
	payV1.UnimplementedPaymentServiceServer
	service service.Service
}

func NewApi(paymentService service.Service) *api {
	return &api{
		service: paymentService,
	}
}
