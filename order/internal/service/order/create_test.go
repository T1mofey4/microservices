package order

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/T1mofey4/microservices/order/internal/model"
)

func (s *ServiceSuite) TestCreateSuccess() {
	var (
		userUUID             = uuid.New()
		partUUID1, partUUID2 = uuid.New(), uuid.New()
		filter               = model.PartsFilter{UUIDs: []uuid.UUID{partUUID1, partUUID2}}
	)

	s.inventoryClient.EXPECT().ListParts(s.ctx, filter).Return([]model.PartInfo{{PartUUID: partUUID1, Price: 0.1}, {PartUUID: partUUID2, Price: 0.2}}, nil)

	var captured model.Order
	s.orderRepository.EXPECT().Create(s.ctx, mock.Anything).Run(func(_ context.Context, order model.Order) {
		captured = order
	}).Return(nil)

	res, err := s.service.Create(s.ctx, userUUID, filter)

	s.Require().NoError(err)
	s.Require().NotEqual(uuid.Nil, captured.OrderUUID)
	s.Require().Equal(userUUID, captured.UserUUID)
	s.Require().Equal(filter.UUIDs, captured.PartUUIDs)
	s.Require().Equal(0.3, captured.TotalPrice)
	s.Require().Equal(model.StatusPendingPayment, captured.Status)
	s.Require().Nil(captured.TransactionUUID)
	s.Require().Nil(captured.PaymentMethod)
	s.Require().Equal(captured, res)
}

func (s *ServiceSuite) TestCreateListPartsError() {
	var (
		userUUID             = uuid.New()
		partUUID1, partUUID2 = uuid.New(), uuid.New()
		someErr              = errors.New("list parts error")
	)

	filter := model.PartsFilter{UUIDs: []uuid.UUID{partUUID1, partUUID2}}

	s.inventoryClient.EXPECT().ListParts(s.ctx, filter).Return(nil, someErr)
	res, err := s.service.Create(s.ctx, userUUID, filter)

	s.Require().Equal(model.Order{}, res)
	s.Require().ErrorIs(err, someErr)
}

func (s *ServiceSuite) TestCreateCreateError() {
	var (
		userUUID             = uuid.New()
		partUUID1, partUUID2 = uuid.New(), uuid.New()
		filter               = model.PartsFilter{UUIDs: []uuid.UUID{partUUID1, partUUID2}}
		someErr              = errors.New("create order error")
	)

	s.inventoryClient.EXPECT().ListParts(s.ctx, filter).Return([]model.PartInfo{{PartUUID: partUUID1, Price: 0.1}, {PartUUID: partUUID2, Price: 0.2}}, nil)

	s.orderRepository.EXPECT().Create(s.ctx, mock.Anything).Return(someErr)
	res, err := s.service.Create(s.ctx, userUUID, filter)

	s.Require().Equal(model.Order{}, res)
	s.Require().ErrorIs(err, someErr)
}
