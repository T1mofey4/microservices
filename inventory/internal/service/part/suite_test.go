package part

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/T1mofey4/microservices/inventory/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	partRepository *mocks.MockPartRepository

	service *partService
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.partRepository = mocks.NewMockPartRepository(s.T())

	s.service = NewService(s.partRepository)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestPartServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
