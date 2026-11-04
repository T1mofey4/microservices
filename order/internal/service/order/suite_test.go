package order

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	clientsMocks "github.com/T1mofey4/microservices/order/internal/client/grpc/mocks"
	"github.com/T1mofey4/microservices/order/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	orderRepository *mocks.MockOrderRepository

	inventoryClient *clientsMocks.MockInventoryClient

	paymentClient *clientsMocks.MockPaymentClient

	service *orderService
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.orderRepository = mocks.NewMockOrderRepository(s.T())
	s.inventoryClient = clientsMocks.NewMockInventoryClient(s.T())
	s.paymentClient = clientsMocks.NewMockPaymentClient(s.T())

	s.service = NewOrderService(s.orderRepository, s.inventoryClient, s.paymentClient)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestOrderServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
