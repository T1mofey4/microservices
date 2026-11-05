package order

import (
	"errors"

	"github.com/google/uuid"

	"github.com/T1mofey4/microservices/order/internal/model"
	repo "github.com/T1mofey4/microservices/order/internal/repository"
)

func (s *ServiceSuite) TestCancelSuccess() {
	var (
		orderUUID  = uuid.New()
		userUUID   = uuid.New()
		partUUIDs  = []uuid.UUID{uuid.New()}
		totalPrice = 100.26
		status     = model.StatusPendingPayment
	)

	expected := model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     status,
	}

	updOrder := model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     model.StatusCancelled,
	}

	s.orderRepository.EXPECT().Get(s.ctx, orderUUID.String()).Return(expected, nil)

	s.orderRepository.EXPECT().Update(s.ctx, orderUUID.String(), model.StatusPendingPayment, updOrder).Return(nil)

	err := s.service.Cancel(s.ctx, orderUUID)

	s.Require().NoError(err)
}

func (s *ServiceSuite) TestCancelGetError() {
	var (
		orderUUID = uuid.New()
		someErr   = errors.New("get order error")
	)

	s.orderRepository.EXPECT().Get(s.ctx, orderUUID.String()).Return(model.Order{}, someErr)

	err := s.service.Cancel(s.ctx, orderUUID)

	s.Require().ErrorIs(err, someErr)
}

func (s *ServiceSuite) TestCancelAlreadyPaid() {
	orderUUID := uuid.New()
	paidOrder := model.Order{
		OrderUUID: orderUUID,
		Status:    model.StatusPaid,
	}

	s.orderRepository.EXPECT().Get(s.ctx, orderUUID.String()).Return(paidOrder, nil)

	err := s.service.Cancel(s.ctx, orderUUID)

	s.Require().ErrorIs(err, model.ErrOrderAlreadyPaid)
}

func (s *ServiceSuite) TestCancelCancelled() {
	orderUUID := uuid.New()
	cancelledOrder := model.Order{
		OrderUUID: orderUUID,
		Status:    model.StatusCancelled,
	}

	s.orderRepository.EXPECT().Get(s.ctx, orderUUID.String()).Return(cancelledOrder, nil)

	err := s.service.Cancel(s.ctx, orderUUID)

	s.Require().ErrorIs(err, model.ErrOrderCancelled)
}

func (s *ServiceSuite) TestCancelInvalidStatus() {
	orderUUID := uuid.New()
	unknownStatusOrder := model.Order{
		OrderUUID: orderUUID,
		Status:    model.OrderStatus("unknown"),
	}

	s.orderRepository.EXPECT().Get(s.ctx, orderUUID.String()).Return(unknownStatusOrder, nil)

	err := s.service.Cancel(s.ctx, orderUUID)

	s.Require().ErrorIs(err, model.ErrInvalidArgument)
}

func (s *ServiceSuite) TestCancelUpdateNotFound() {
	var (
		orderUUID  = uuid.New()
		userUUID   = uuid.New()
		partUUIDs  = []uuid.UUID{uuid.New()}
		totalPrice = 100.26
		status     = model.StatusPendingPayment
	)

	expected := model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     status,
	}

	updOrder := model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     model.StatusCancelled,
	}

	s.orderRepository.EXPECT().Get(s.ctx, orderUUID.String()).Return(expected, nil)

	s.orderRepository.EXPECT().Update(s.ctx, orderUUID.String(), model.StatusPendingPayment, updOrder).Return(model.ErrOrderNotFound)

	err := s.service.Cancel(s.ctx, orderUUID)

	s.Require().ErrorIs(err, model.ErrOrderNotFound)
}

func (s *ServiceSuite) TestCancelUpdateOtherError() {
	var (
		orderUUID  = uuid.New()
		userUUID   = uuid.New()
		partUUIDs  = []uuid.UUID{uuid.New()}
		totalPrice = 100.26
		status     = model.StatusPendingPayment
		someErr    = errors.New("update error")
	)

	expected := model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     status,
	}

	updOrder := model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     model.StatusCancelled,
	}

	s.orderRepository.EXPECT().Get(s.ctx, orderUUID.String()).Return(expected, nil)

	s.orderRepository.EXPECT().Update(s.ctx, orderUUID.String(), model.StatusPendingPayment, updOrder).Return(someErr)

	err := s.service.Cancel(s.ctx, orderUUID)

	s.Require().ErrorIs(err, someErr)
}

func (s *ServiceSuite) TestCancelUpdateStatusConflictPaid() {
	var (
		orderUUID  = uuid.New()
		userUUID   = uuid.New()
		partUUIDs  = []uuid.UUID{uuid.New()}
		totalPrice = 100.26
		status     = model.StatusPendingPayment
	)

	pendingOrder := model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     status,
	}

	retryOrder := model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     model.StatusPaid,
	}

	updOrder := model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     model.StatusCancelled,
	}

	s.orderRepository.EXPECT().Get(s.ctx, orderUUID.String()).Return(pendingOrder, nil).Once()

	s.orderRepository.EXPECT().Get(s.ctx, orderUUID.String()).Return(retryOrder, nil).Once()

	s.orderRepository.EXPECT().Update(s.ctx, orderUUID.String(), model.StatusPendingPayment, updOrder).Return(repo.ErrStatusConflict)

	err := s.service.Cancel(s.ctx, orderUUID)

	s.Require().ErrorIs(err, model.ErrOrderAlreadyPaid)
}

func (s *ServiceSuite) TestCancelUpdateStatusConflictCancelled() {
	var (
		orderUUID  = uuid.New()
		userUUID   = uuid.New()
		partUUIDs  = []uuid.UUID{uuid.New()}
		totalPrice = 100.26
		status     = model.StatusPendingPayment
	)

	pendingOrder := model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     status,
	}

	retryOrder := model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     model.StatusCancelled,
	}

	updOrder := model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     model.StatusCancelled,
	}

	s.orderRepository.EXPECT().Get(s.ctx, orderUUID.String()).Return(pendingOrder, nil).Once()

	s.orderRepository.EXPECT().Get(s.ctx, orderUUID.String()).Return(retryOrder, nil).Once()

	s.orderRepository.EXPECT().Update(s.ctx, orderUUID.String(), model.StatusPendingPayment, updOrder).Return(repo.ErrStatusConflict)

	err := s.service.Cancel(s.ctx, orderUUID)

	s.Require().ErrorIs(err, model.ErrOrderCancelled)
}

func (s *ServiceSuite) TestCancelUpdateStatusConflictRetryUnknown() {
	var (
		orderUUID  = uuid.New()
		userUUID   = uuid.New()
		partUUIDs  = []uuid.UUID{uuid.New()}
		totalPrice = 100.26
		status     = model.StatusPendingPayment
	)

	pendingOrder := model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     status,
	}

	retryOrder := model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     model.OrderStatus("unknown"),
	}

	updOrder := model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     model.StatusCancelled,
	}

	s.orderRepository.EXPECT().Get(s.ctx, orderUUID.String()).Return(pendingOrder, nil).Once()

	s.orderRepository.EXPECT().Get(s.ctx, orderUUID.String()).Return(retryOrder, nil).Once()

	s.orderRepository.EXPECT().Update(s.ctx, orderUUID.String(), model.StatusPendingPayment, updOrder).Return(repo.ErrStatusConflict)

	err := s.service.Cancel(s.ctx, orderUUID)

	s.Require().ErrorIs(err, repo.ErrStatusConflict)
}

func (s *ServiceSuite) TestCancelUpdateStatusConflictRetryGetError() {
	var (
		orderUUID  = uuid.New()
		userUUID   = uuid.New()
		partUUIDs  = []uuid.UUID{uuid.New()}
		totalPrice = 100.26
		status     = model.StatusPendingPayment
		retryErr   = errors.New("retry error")
	)

	pendingOrder := model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     status,
	}

	updOrder := model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     model.StatusCancelled,
	}

	s.orderRepository.EXPECT().Get(s.ctx, orderUUID.String()).Return(pendingOrder, nil).Once()

	s.orderRepository.EXPECT().Get(s.ctx, orderUUID.String()).Return(model.Order{}, retryErr).Once()

	s.orderRepository.EXPECT().Update(s.ctx, orderUUID.String(), model.StatusPendingPayment, updOrder).Return(repo.ErrStatusConflict)

	err := s.service.Cancel(s.ctx, orderUUID)

	s.Require().ErrorIs(err, retryErr)
}
