package part

import (
	"time"

	repoModel "github.com/T1mofey4/microservices/inventory/internal/repository/model"
)

func initParts() []repoModel.Part {
	parts := []repoModel.Part{
		{UUID: "11111111-1111-1111-1111-111111111111",
			Name:          "Main Engine",
			Description:   "The primary engine unit for the main vehicle system.",
			Price:         111.00,
			StockQuantity: 111,
			Category:      repoModel.CategoryEngine,
			Dimensions: repoModel.Dimensions{
				Length: 0.1,
				Width:  0.2,
				Height: 0.3,
				Weight: 0.5,
			},
			Manufacturer: repoModel.Manufacturer{
				Name:    "Test Manufacturer Corp",
				Country: "USA",
				Website: "http://testcorp.com",
			},
			Tags:      []string{"electronics", "tools"},
			CreatedAt: time.Now(),
		},
		{UUID: "22222222-2222-2222-2222-222222222222",
			Name:          "Main Fuel Tank",
			Description:   "Holds the primary fuel source for operation.",
			Price:         55.50,
			StockQuantity: 500,
			Category:      repoModel.CategoryFuel,
			Dimensions: repoModel.Dimensions{
				Length: 1.0,
				Width:  0.5,
				Height: 0.2,
				Weight: 1.2,
			},
			Manufacturer: repoModel.Manufacturer{
				Name:    "WoodWorks Inc",
				Country: "CAN",
				Website: "http://woodworks.ca",
			},
			Tags:      []string{"wood", "hardware"},
			CreatedAt: time.Now(),
		},
		{UUID: "33333333-3333-3333-3333-333333333333",
			Name:          "Part3",
			Description:   "High-performance component for advanced machinery.",
			Price:         1200.00,
			StockQuantity: 10,
			Category:      repoModel.CategoryEngine,
			Dimensions: repoModel.Dimensions{
				Length: 2.0,
				Width:  1.5,
				Height: 0.5,
				Weight: 25.0,
			},
			Manufacturer: repoModel.Manufacturer{
				Name:    "ElectroTech Global",
				Country: "JPN",
				Website: "http://electrotech.jp",
			},
			Tags:      []string{"electronics", "high-end"},
			CreatedAt: time.Now(),
		},
		{UUID: "44444444-4444-4444-4444-444444444444",
			Name:          "Part4",
			Description:   "Small, universal fastening piece.",
			Price:         12.99,
			StockQuantity: 2000,
			Category:      repoModel.CategoryFuel,
			Dimensions: repoModel.Dimensions{
				Length: 0.1,
				Width:  0.1,
				Height: 0.1,
				Weight: 0.05,
			},
			Manufacturer: repoModel.Manufacturer{
				Name:    "Fasteners Ltd",
				Country: "USA",
				Website: "http://fasteners.com",
			},
			Tags:      []string{"fasteners", "small"},
			CreatedAt: time.Now(),
		},
	}
	return parts
}
