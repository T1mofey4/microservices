package payment

import def "github.com/T1mofey4/microservices/payment/internal/service"

var _ def.Service = (*service)(nil)

type service struct{}

func NewService() *service {
	return &service{}
}
