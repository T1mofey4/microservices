package part

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/T1mofey4/microservices/inventory/internal/model"
)

func (s *ServiceSuite) TestGetPartSuccess() {
	var (
		partUuid      = uuid.New()
		name          = "Main Engine"
		description   = "The primary engine unit for the main vehicle system."
		price         = 111.00
		stockQuantity = 111
		category      = model.CategoryEngine
		dimensions    = model.Dimensions{
			Length: 0.1,
			Width:  0.2,
			Height: 0.3,
			Weight: 0.5,
		}
		manufacturer = model.Manufacturer{
			Name:    "Test Manufacturer Corp",
			Country: "USA",
			Website: "http://testcorp.com",
		}
		tags      = []string{"electronics", "tools"}
		createdAt = time.Now()

		part = model.Part{
			UUID:          partUuid.String(),
			Name:          name,
			Description:   description,
			Price:         price,
			StockQuantity: int64(stockQuantity),
			Category:      category,
			Dimensions:    dimensions,
			Manufacturer:  manufacturer,
			Tags:          tags,
			CreatedAt:     createdAt,
		}
	)

	s.partRepository.EXPECT().Get(s.ctx, partUuid.String()).Return(&part, nil)

	res, err := s.service.GetPart(s.ctx, partUuid.String())
	s.Require().NoError(err)
	s.Require().Equal(&part, res)
}

func (s *ServiceSuite) TestGetPartNotFound() {
	partUuid := uuid.New()

	s.partRepository.EXPECT().Get(s.ctx, partUuid.String()).Return(nil, model.ErrPartNotFound)

	res, err := s.service.GetPart(s.ctx, partUuid.String())
	s.Require().Nil(res)
	s.Require().ErrorIs(err, model.ErrPartNotFound)
}

func (s *ServiceSuite) TestGetPartRepositoryError() {
	partUuid := uuid.New()
	someErr := errors.New("repository error")

	s.partRepository.EXPECT().Get(s.ctx, partUuid.String()).Return(nil, someErr)

	res, err := s.service.GetPart(s.ctx, partUuid.String())
	s.Require().Nil(res)
	s.Require().ErrorIs(err, someErr)
	s.Require().NotErrorIs(err, model.ErrPartNotFound)
}
