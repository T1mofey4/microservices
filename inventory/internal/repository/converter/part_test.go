package converter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/T1mofey4/microservices/inventory/internal/model"
	repoModel "github.com/T1mofey4/microservices/inventory/internal/repository/model"
)

func TestPartToProto(t *testing.T) {
	testTime := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name         string
		serviceModel repoModel.Part
		expected     *model.Part
	}{
		{
			name: "Valid test data",
			serviceModel: repoModel.Part{
				UUID:          "11111111-1111-1111-1111-111111111111",
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
				CreatedAt: testTime,
				UpdatedAt: &testTime,
			},
			expected: &model.Part{
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
				CreatedAt: testTime,
				UpdatedAt: &testTime,
			},
		},
		{
			name: "Nil UpdatedAt and empty tags",
			serviceModel: repoModel.Part{
				UUID:      "22222222-2222-2222-2222-222222222222",
				CreatedAt: testTime,
				UpdatedAt: nil,
			},

			expected: &model.Part{
				UUID:         "22222222-2222-2222-2222-222222222222",
				Dimensions:   model.Dimensions{},
				Manufacturer: model.Manufacturer{},
				CreatedAt:    testTime,
				UpdatedAt:    nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := PartToModel(tt.serviceModel)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestCategoryToRepoModel(t *testing.T) {
	tests := []struct {
		name     string
		category model.Category
		expected repoModel.Category
	}{
		{name: "Category Engine", category: model.CategoryEngine, expected: repoModel.CategoryEngine},
		{name: "Category Fuel", category: model.CategoryFuel, expected: repoModel.CategoryFuel},
		{name: "Category Porthole", category: model.CategoryPorthole, expected: repoModel.CategoryPorthole},
		{name: "Category Wing", category: model.CategoryWing, expected: repoModel.CategoryWing},
		{name: "Category Unspecified", category: model.CategoryUnspecified, expected: repoModel.CategoryUnspecified},
		{name: "Wrong Category", category: model.Category(111), expected: repoModel.CategoryUnspecified},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := CategoryToRepoModel(tt.category)
			require.Equal(t, tt.expected, result)
		})
	}
}
