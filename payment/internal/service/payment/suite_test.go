package payment

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
	suite.Suite
	//nolint:containedctx
	ctx context.Context

	service *paymentService
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.service = NewPaymentService()
}

func (s *ServiceSuite) TearDownTest() {
}

func TestPaymentServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
