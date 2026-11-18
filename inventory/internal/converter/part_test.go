package converter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/T1mofey4/microservices/inventory/internal/model"
	inventoryV1 "github.com/T1mofey4/microservices/shared/pkg/proto/inventory/v1"
)

func TestPartToProto(t *testing.T) {
	testTime := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name         string
		serviceModel *model.Part
		expected     *inventoryV1.Part
	}{
		{
			name: "Valid test data",
			serviceModel: &model.Part{
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
			expected: &inventoryV1.Part{
				Uuid:          "11111111-1111-1111-1111-111111111111",
				Name:          "Main Engine",
				Description:   "The primary engine unit for the main vehicle system.",
				Price:         111.00,
				StockQuantity: 111,
				Category:      inventoryV1.Category_CATEGORY_ENGINE,
				Dimensions: &inventoryV1.Dimensions{
					Length: 0.1,
					Width:  0.2,
					Height: 0.3,
					Weight: 0.5,
				},
				Manufacturer: &inventoryV1.Manufacturer{
					Name:    "Test Manufacturer Corp",
					Country: "USA",
					Website: "http://testcorp.com",
				},
				Tags:      []string{"electronics", "tools"},
				CreatedAt: timestamppb.New(testTime),
				UpdatedAt: timestamppb.New(testTime),
			},
		},
		{
			name: "Nil UpdatedAt and empty tags",
			serviceModel: &model.Part{
				UUID:      "22222222-2222-2222-2222-222222222222",
				CreatedAt: testTime,
				UpdatedAt: nil,
			},

			expected: &inventoryV1.Part{
				Uuid:         "22222222-2222-2222-2222-222222222222",
				Dimensions:   &inventoryV1.Dimensions{},
				Manufacturer: &inventoryV1.Manufacturer{},
				CreatedAt:    timestamppb.New(testTime),
				UpdatedAt:    nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := PartToProto(tt.serviceModel)
			require.True(t, proto.Equal(tt.expected, result))
		})
	}
}

func TestFilterToModel(t *testing.T) {
	tests := []struct {
		name       string
		protoModel *inventoryV1.PartsFilter
		expected   model.PartsFilter
	}{
		{
			name: "Valid data",
			protoModel: &inventoryV1.PartsFilter{
				Uuids: []string{
					"11111111-1111-1111-1111-111111111111", "22222222-2222-2222-2222-222222222222",
				},
				Names:                 []string{"Main Engine"},
				Categories:            []inventoryV1.Category{inventoryV1.Category_CATEGORY_ENGINE},
				ManufacturerCountries: []string{"USA", "CAN"},
			},

			expected: model.PartsFilter{
				UUIDs: []string{
					"11111111-1111-1111-1111-111111111111", "22222222-2222-2222-2222-222222222222",
				},
				Names:                 []string{"Main Engine"},
				Categories:            []model.Category{model.CategoryEngine},
				ManufacturerCountries: []string{"USA", "CAN"},
			},
		},

		{
			name:       "Empty filter",
			protoModel: nil,
			expected:   model.PartsFilter{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := FilterToModel(tt.protoModel)
			require.Equal(t, tt.expected, result)
		})
	}
}
