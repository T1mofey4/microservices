package v1

import (
	orderV1 "github.com/T1mofey4/microservices/shared/pkg/openapi/order/v1"

	def "github.com/T1mofey4/microservices/order/internal/service"
)

var _ orderV1.Handler = (*api)(nil)

type api struct {
	orderService def.OrderService
}

func NewApi(orderService def.OrderService) *api {
	return &api{
		orderService: orderService,
	}
}
