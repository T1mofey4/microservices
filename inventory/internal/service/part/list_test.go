package part

import (
	"errors"
	"time"

	"github.com/T1mofey4/microservices/inventory/internal/model"
)

func (s *ServiceSuite) TestListPartsSuccess() {
	filter := model.PartsFilter{
		UUIDs: []string{
			"22222222-2222-2222-2222-222222222222",
			"11111111-1111-1111-1111-111111111111",
		},
		Names:                 []string{"Main Fuel Tank", "Main Engine"},
		Categories:            []model.Category{model.CategoryEngine, model.CategoryFuel},
		ManufacturerCountries: []string{"USA", "JPN"},
		Tags:                  []string{"electronics"},
	}

	parts := []model.Part{
		{
			UUID:          "11111111-1111-1111-1111-111111111111",
			Name:          "Main Engine",
			Description:   "The primary engine unit for the main vehicle system.",
			Price:         111.00,
			StockQuantity: 111,
			Category:      model.CategoryEngine,
			Dimensions: model.Dimensions{
				Length: 0.1,
				Width:  0.2,
				Height: 0.3,
				Weight: 0.5,
			},
			Manufacturer: model.Manufacturer{
				Name:    "Test Manufacturer Corp",
				Country: "USA",
				Website: "http://testcorp.com",
			},
			Tags:      []string{"electronics", "tools"},
			CreatedAt: time.Now(),
		},
		{
			UUID:          "22222222-2222-2222-2222-222222222222",
			Name:          "Main Fuel Tank",
			Description:   "Holds the primary fuel source for operation.",
			Price:         55.50,
			StockQuantity: 500,
			Category:      model.CategoryFuel,
			Dimensions: model.Dimensions{
				Length: 1.0,
				Width:  0.5,
				Height: 0.2,
				Weight: 1.2,
			},
			Manufacturer: model.Manufacturer{
				Name:    "WoodWorks Inc",
				Country: "CAN",
				Website: "http://woodworks.ca",
			},
			Tags:      []string{"wood", "hardware"},
			CreatedAt: time.Now(),
		},
	}

	s.partRepository.EXPECT().ListParts(s.ctx, filter).Return(parts, nil)

	res, err := s.service.ListParts(s.ctx, filter)
	s.Require().NoError(err)
	s.Require().Equal(parts, res)
}

func (s *ServiceSuite) TestListPartsEmptyResult() {
	filter := model.PartsFilter{}
	parts := []model.Part{}

	s.partRepository.EXPECT().ListParts(s.ctx, filter).Return(parts, nil)

	res, err := s.service.ListParts(s.ctx, filter)
	s.Require().NoError(err)
	s.Require().Equal(parts, res)
}

func (s *ServiceSuite) TestListPartsRepositoryError() {
	filter := model.PartsFilter{
		UUIDs: []string{
			"22222222-2222-2222-2222-222222222222",
			"11111111-1111-1111-1111-111111111111",
		},
		Names:                 []string{"Main Fuel Tank", "Main Engine"},
		Categories:            []model.Category{model.CategoryEngine, model.CategoryFuel},
		ManufacturerCountries: []string{"USA", "JPN"},
		Tags:                  []string{"electronics"},
	}

	someErr := errors.New("repository error")

	s.partRepository.EXPECT().ListParts(s.ctx, filter).Return(nil, someErr)

	res, err := s.service.ListParts(s.ctx, filter)
	s.Require().Nil(res)
	s.Require().ErrorIs(err, someErr)
}
