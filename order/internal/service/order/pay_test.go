package order

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/T1mofey4/microservices/order/internal/model"
	repo "github.com/T1mofey4/microservices/order/internal/repository"
)

func (s *ServiceSuite) TestPaySuccess() {
	var (
		orderUUID       = uuid.New()
		userUUID        = uuid.New()
		partUUIDs       = []uuid.UUID{uuid.New()}
		totalPrice      = 100.26
		status          = model.StatusPendingPayment
		paymentMethod   = model.PaymentMethodCreditCard
		transactionUUID = uuid.New()
	)

	expected := model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     status,
	}

	s.orderRepository.EXPECT().Get(s.ctx, orderUUID.String()).Return(expected, nil)

	s.paymentClient.EXPECT().PayOrder(s.ctx, userUUID, orderUUID, paymentMethod).Return(transactionUUID, nil)

	var captured model.Order
	s.orderRepository.EXPECT().
		Update(s.ctx, orderUUID.String(), model.StatusPendingPayment, mock.Anything).
		Run(func(_ context.Context, _ string, _ model.OrderStatus, upd model.Order) {
			captured = upd
		}).
		Return(nil)

	res, err := s.service.Pay(s.ctx, orderUUID, paymentMethod)

	s.Require().NoError(err)
	s.Require().Equal(transactionUUID, res)
	s.Require().Equal(orderUUID, captured.OrderUUID)
	s.Require().Equal(userUUID, captured.UserUUID)
	s.Require().Equal(partUUIDs, captured.PartUUIDs)
	s.Require().Equal(totalPrice, captured.TotalPrice)
	s.Require().Equal(&transactionUUID, captured.TransactionUUID)
	s.Require().Equal(&paymentMethod, captured.PaymentMethod)
	s.Require().Equal(model.StatusPaid, captured.Status)
}

func (s *ServiceSuite) TestPayGetError() {
	someErr := errors.New("TestPayGetError")
	orderUUID := uuid.New()
	paymentMethod := model.PaymentMethodCreditCard

	s.orderRepository.EXPECT().Get(s.ctx, orderUUID.String()).Return(model.Order{}, someErr)

	res, err := s.service.Pay(s.ctx, orderUUID, paymentMethod)
	s.Require().Equal(uuid.Nil, res)
	s.Require().ErrorIs(err, someErr)
}

func (s *ServiceSuite) TestPayAlreadyPaid() {
	var (
		orderUUID     = uuid.New()
		userUUID      = uuid.New()
		partUUIDs     = []uuid.UUID{uuid.New()}
		totalPrice    = 100.26
		status        = model.StatusPaid
		paymentMethod = model.PaymentMethodCreditCard
	)

	expected := model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     status,
	}

	s.orderRepository.EXPECT().Get(s.ctx, orderUUID.String()).Return(expected, nil)

	res, err := s.service.Pay(s.ctx, orderUUID, paymentMethod)
	s.Require().Equal(uuid.Nil, res)
	s.Require().ErrorIs(err, model.ErrOrderAlreadyPaid)
}

func (s *ServiceSuite) TestPayCancelled() {
	var (
		orderUUID     = uuid.New()
		userUUID      = uuid.New()
		partUUIDs     = []uuid.UUID{uuid.New()}
		totalPrice    = 100.26
		status        = model.StatusCancelled
		paymentMethod = model.PaymentMethodCreditCard
	)

	expected := model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     status,
	}

	s.orderRepository.EXPECT().Get(s.ctx, orderUUID.String()).Return(expected, nil)

	res, err := s.service.Pay(s.ctx, orderUUID, paymentMethod)
	s.Require().Equal(uuid.Nil, res)
	s.Require().ErrorIs(err, model.ErrOrderCancelled)
}

func (s *ServiceSuite) TestPayInvalidStatus() {
	var (
		orderUUID     = uuid.New()
		userUUID      = uuid.New()
		partUUIDs     = []uuid.UUID{uuid.New()}
		totalPrice    = 100.26
		status        = model.OrderStatus("unknown")
		paymentMethod = model.PaymentMethodCreditCard
	)

	expected := model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     status,
	}

	s.orderRepository.EXPECT().Get(s.ctx, orderUUID.String()).Return(expected, nil)

	res, err := s.service.Pay(s.ctx, orderUUID, paymentMethod)
	s.Require().Equal(uuid.Nil, res)
	s.Require().ErrorIs(err, model.ErrInvalidArgument)
}

func (s *ServiceSuite) TestPayPaymentError() {
	var (
		someErr       = errors.New("payment error")
		orderUUID     = uuid.New()
		userUUID      = uuid.New()
		partUUIDs     = []uuid.UUID{uuid.New()}
		totalPrice    = 100.26
		status        = model.StatusPendingPayment
		paymentMethod = model.PaymentMethodCreditCard
	)

	expected := model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     status,
	}

	s.orderRepository.EXPECT().Get(s.ctx, orderUUID.String()).Return(expected, nil)

	s.paymentClient.EXPECT().PayOrder(s.ctx, userUUID, orderUUID, paymentMethod).Return(uuid.Nil, someErr)
	res, err := s.service.Pay(s.ctx, orderUUID, paymentMethod)

	s.Require().Equal(uuid.Nil, res)
	s.Require().ErrorIs(err, someErr)
}

func (s *ServiceSuite) TestPayUpdateNotFound() {
	var (
		orderUUID       = uuid.New()
		userUUID        = uuid.New()
		partUUIDs       = []uuid.UUID{uuid.New()}
		totalPrice      = 100.26
		status          = model.StatusPendingPayment
		paymentMethod   = model.PaymentMethodCreditCard
		transactionUUID = uuid.New()
	)

	expected := model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     status,
	}

	s.orderRepository.EXPECT().Get(s.ctx, orderUUID.String()).Return(expected, nil)

	s.paymentClient.EXPECT().PayOrder(s.ctx, userUUID, orderUUID, paymentMethod).Return(transactionUUID, nil)

	s.orderRepository.EXPECT().
		Update(s.ctx, orderUUID.String(), model.StatusPendingPayment, mock.Anything).
		Return(model.ErrOrderNotFound)

	res, err := s.service.Pay(s.ctx, orderUUID, paymentMethod)

	s.Require().Equal(uuid.Nil, res)
	s.Require().ErrorIs(err, model.ErrOrderNotFound)
}

func (s *ServiceSuite) TestPayUpdateOtherError() {
	var (
		someErr         = errors.New("update error")
		orderUUID       = uuid.New()
		userUUID        = uuid.New()
		partUUIDs       = []uuid.UUID{uuid.New()}
		totalPrice      = 100.26
		status          = model.StatusPendingPayment
		paymentMethod   = model.PaymentMethodCreditCard
		transactionUUID = uuid.New()
	)

	expected := model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     status,
	}

	s.orderRepository.EXPECT().Get(s.ctx, orderUUID.String()).Return(expected, nil)

	s.paymentClient.EXPECT().PayOrder(s.ctx, userUUID, orderUUID, paymentMethod).Return(transactionUUID, nil)

	s.orderRepository.EXPECT().
		Update(s.ctx, orderUUID.String(), model.StatusPendingPayment, mock.Anything).
		Return(someErr)

	res, err := s.service.Pay(s.ctx, orderUUID, paymentMethod)

	s.Require().Equal(uuid.Nil, res)
	s.Require().ErrorIs(err, someErr)
}

func (s *ServiceSuite) TestPayUpdateStatusConflictPaid() {
	var (
		orderUUID       = uuid.New()
		userUUID        = uuid.New()
		partUUIDs       = []uuid.UUID{uuid.New()}
		totalPrice      = 100.26
		paymentMethod   = model.PaymentMethodCreditCard
		transactionUUID = uuid.New()
	)

	pendingOrder := model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     model.StatusPendingPayment,
	}

	paidOrder := model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     model.StatusPaid,
	}

	s.orderRepository.EXPECT().Get(s.ctx, orderUUID.String()).Return(pendingOrder, nil).Once()

	s.orderRepository.EXPECT().Get(s.ctx, orderUUID.String()).Return(paidOrder, nil).Once()

	s.paymentClient.EXPECT().PayOrder(s.ctx, userUUID, orderUUID, paymentMethod).Return(transactionUUID, nil)

	s.orderRepository.EXPECT().
		Update(s.ctx, orderUUID.String(), model.StatusPendingPayment, mock.Anything).
		Return(repo.ErrStatusConflict)

	res, err := s.service.Pay(s.ctx, orderUUID, paymentMethod)
	s.Require().Equal(uuid.Nil, res)
	s.Require().ErrorIs(err, model.ErrOrderAlreadyPaid)
}

func (s *ServiceSuite) TestPayUpdateStatusConflictCancelled() {
	var (
		orderUUID       = uuid.New()
		userUUID        = uuid.New()
		partUUIDs       = []uuid.UUID{uuid.New()}
		totalPrice      = 100.26
		paymentMethod   = model.PaymentMethodCreditCard
		transactionUUID = uuid.New()
	)

	pendingOrder := model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     model.StatusPendingPayment,
	}

	cancelledOrder := model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     model.StatusCancelled,
	}

	s.orderRepository.EXPECT().Get(s.ctx, orderUUID.String()).Return(pendingOrder, nil).Once()

	s.orderRepository.EXPECT().Get(s.ctx, orderUUID.String()).Return(cancelledOrder, nil).Once()

	s.paymentClient.EXPECT().PayOrder(s.ctx, userUUID, orderUUID, paymentMethod).Return(transactionUUID, nil)

	s.orderRepository.EXPECT().
		Update(s.ctx, orderUUID.String(), model.StatusPendingPayment, mock.Anything).
		Return(repo.ErrStatusConflict)

	res, err := s.service.Pay(s.ctx, orderUUID, paymentMethod)
	s.Require().Equal(uuid.Nil, res)
	s.Require().ErrorIs(err, model.ErrOrderCancelled)
}

func (s *ServiceSuite) TestPayUpdateStatusConflictRetryUnknown() {
	var (
		orderUUID       = uuid.New()
		userUUID        = uuid.New()
		partUUIDs       = []uuid.UUID{uuid.New()}
		totalPrice      = 100.26
		paymentMethod   = model.PaymentMethodCreditCard
		transactionUUID = uuid.New()
	)

	pendingOrder := model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     model.StatusPendingPayment,
	}

	statusConflictOrder := model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     model.OrderStatus("unknown"),
	}

	s.orderRepository.EXPECT().Get(s.ctx, orderUUID.String()).Return(pendingOrder, nil).Once()

	s.orderRepository.EXPECT().Get(s.ctx, orderUUID.String()).Return(statusConflictOrder, nil).Once()

	s.paymentClient.EXPECT().PayOrder(s.ctx, userUUID, orderUUID, paymentMethod).Return(transactionUUID, nil)

	s.orderRepository.EXPECT().
		Update(s.ctx, orderUUID.String(), model.StatusPendingPayment, mock.Anything).
		Return(repo.ErrStatusConflict)

	res, err := s.service.Pay(s.ctx, orderUUID, paymentMethod)
	s.Require().Equal(uuid.Nil, res)
	s.Require().ErrorIs(err, repo.ErrStatusConflict)
}

func (s *ServiceSuite) TestPayUpdateStatusConflictRetryGetError() {
	var (
		orderUUID       = uuid.New()
		userUUID        = uuid.New()
		partUUIDs       = []uuid.UUID{uuid.New()}
		totalPrice      = 100.26
		paymentMethod   = model.PaymentMethodCreditCard
		transactionUUID = uuid.New()
		retryErr        = errors.New("retry get failed")
	)

	pendingOrder := model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     model.StatusPendingPayment,
	}

	s.orderRepository.EXPECT().Get(s.ctx, orderUUID.String()).Return(pendingOrder, nil).Once()

	s.orderRepository.EXPECT().Get(s.ctx, orderUUID.String()).Return(model.Order{}, retryErr).Once()

	s.paymentClient.EXPECT().PayOrder(s.ctx, userUUID, orderUUID, paymentMethod).Return(transactionUUID, nil)

	s.orderRepository.EXPECT().
		Update(s.ctx, orderUUID.String(), model.StatusPendingPayment, mock.Anything).
		Return(repo.ErrStatusConflict)

	res, err := s.service.Pay(s.ctx, orderUUID, paymentMethod)
	s.Require().Equal(uuid.Nil, res)
	s.Require().ErrorIs(err, retryErr)
}
