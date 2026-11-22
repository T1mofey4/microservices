package app

import (
	apiV1 "github.com/T1mofey4/microservices/payment/internal/api/payment/v1"
	"github.com/T1mofey4/microservices/payment/internal/service"
	payV1 "github.com/T1mofey4/microservices/payment/internal/service/payment"
	paymentV1 "github.com/T1mofey4/microservices/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	paymentV1API   paymentV1.PaymentServiceServer
	paymentService service.PaymentService
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) PaymentV1API() paymentV1.PaymentServiceServer {
	if d.paymentV1API == nil {
		d.paymentV1API = apiV1.NewApi(d.PaymentService())
	}

	return d.paymentV1API
}

func (d *diContainer) PaymentService() service.PaymentService {
	if d.paymentService == nil {
		d.paymentService = payV1.NewPaymentService()
	}

	return d.paymentService
}
