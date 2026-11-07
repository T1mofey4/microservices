package order

import (
	"errors"

	"github.com/google/uuid"

	"github.com/T1mofey4/microservices/order/internal/model"
)

func (s *ServiceSuite) TestGetSuccess() {
	var (
		orderUUID  = uuid.New()
		userUUID   = uuid.New()
		partUUIDs  = []uuid.UUID{uuid.New()}
		totalPrice = 100.26
		status     = model.StatusPendingPayment

		expected = model.Order{
			OrderUUID:  orderUUID,
			UserUUID:   userUUID,
			PartUUIDs:  partUUIDs,
			TotalPrice: totalPrice,
			Status:     status,
		}
	)

	s.orderRepository.EXPECT().Get(s.ctx, orderUUID.String()).Return(expected, nil)

	res, err := s.service.Get(s.ctx, orderUUID)
	s.Require().NoError(err)
	s.Require().Equal(expected, res)
}

func (s *ServiceSuite) TestGetNotFound() {
	orderUUID := uuid.New()

	s.orderRepository.EXPECT().Get(s.ctx, orderUUID.String()).Return(model.Order{}, model.ErrOrderNotFound)

	_, err := s.service.Get(s.ctx, orderUUID)
	s.Require().ErrorIs(err, model.ErrOrderNotFound)
}

func (s *ServiceSuite) TestGetRepositoryError() {
	orderUUID := uuid.New()
	someErr := errors.New("error")

	s.orderRepository.EXPECT().Get(s.ctx, orderUUID.String()).Return(model.Order{}, someErr)

	_, err := s.service.Get(s.ctx, orderUUID)
	s.Require().ErrorIs(err, someErr)
	s.Require().NotErrorIs(err, model.ErrOrderNotFound)
}
