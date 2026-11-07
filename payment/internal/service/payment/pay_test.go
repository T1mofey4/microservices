package payment

import (
	"github.com/google/uuid"

	"github.com/T1mofey4/microservices/payment/internal/model"
)

func (s *ServiceSuite) TestPayOrderSuccess() {
	var (
		orderUUID = uuid.New().String()
		userUUID  = uuid.New().String()
		method    = model.PaymentMethodCreditCard
	)

	res, err := s.service.PayOrder(s.ctx, orderUUID, userUUID, method)
	s.Require().NoError(err)
	s.Require().NotEqual("", res)
	// проверка на валидный uuid
	_, parseErr := uuid.Parse(res)
	s.Require().NoError(parseErr)
}

func (s *ServiceSuite) TestPayOrderInvalidMethod() {
	var (
		orderUUID = uuid.New().String()
		userUUID  = uuid.New().String()
		method    = model.PaymentMethodUnspecified
	)

	res, err := s.service.PayOrder(s.ctx, orderUUID, userUUID, method)
	s.Require().Equal("", res)
	s.Require().ErrorIs(err, model.ErrInvalidPaymentMethod)
}
