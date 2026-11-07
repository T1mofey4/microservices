package part

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/T1mofey4/microservices/inventory/internal/model"
	repoModel "github.com/T1mofey4/microservices/inventory/internal/repository/model"
)

func TestFilterByUUID(t *testing.T) {
	parts := []repoModel.Part{
		{UUID: "11111111-1111-1111-1111-111111111111"},
		{UUID: "22222222-2222-2222-2222-222222222222"},
		{UUID: "33333333-3333-3333-3333-333333333333"},
		{UUID: "44444444-4444-4444-4444-444444444444"},
	}

	tests := []struct {
		name     string
		filter   []string
		parts    []repoModel.Part
		expected []repoModel.Part
	}{
		{name: "Empty filter", filter: nil, parts: parts, expected: parts},
		{name: "One uuid filter", filter: []string{parts[0].UUID}, parts: parts, expected: []repoModel.Part{parts[0]}},
		{name: "Two uuid filter", filter: []string{parts[0].UUID, parts[1].UUID}, parts: parts, expected: []repoModel.Part{parts[0], parts[1]}},
		{name: "Non-existing filter", filter: []string{"f47ac10b-58cc-4372-a567-0e02b2c3d479"}, parts: parts, expected: []repoModel.Part{}},
		{name: "One exist and one non-existing", filter: []string{"f47ac10b-58cc-4372-a567-0e02b2c3d479", parts[3].UUID}, parts: parts, expected: []repoModel.Part{parts[3]}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := filterByUUID(tt.filter, tt.parts)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestFilterByName(t *testing.T) {
	parts := []repoModel.Part{
		{Name: "Main engine"},
		{Name: "Fuel tank"},
		{Name: "Left wing"},
	}

	tests := []struct {
		name     string
		filter   []string
		parts    []repoModel.Part
		expected []repoModel.Part
	}{
		{name: "Empty filter", filter: nil, parts: parts, expected: parts},
		{name: "Filter by one name", filter: []string{parts[0].Name}, parts: parts, expected: []repoModel.Part{parts[0]}},
		{name: "Filter by two names", filter: []string{parts[0].Name, parts[1].Name}, parts: parts, expected: []repoModel.Part{parts[0], parts[1]}},
		{name: "Filter by non-existing name", filter: []string{"Engine"}, parts: parts, expected: []repoModel.Part{}},
		{name: "Filter by one exist and one non-existing", filter: []string{parts[2].Name, "Wheel"}, parts: parts, expected: []repoModel.Part{parts[2]}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := filterByName(tt.filter, tt.parts)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestFilterByCategory(t *testing.T) {
	parts := []repoModel.Part{
		{Name: "Main engine", Category: repoModel.CategoryEngine},
		{Name: "Fuel tank", Category: repoModel.CategoryFuel},
		{Name: "Left wing", Category: repoModel.CategoryWing},
	}

	tests := []struct {
		name     string
		filter   []model.Category
		parts    []repoModel.Part
		expected []repoModel.Part
	}{
		{name: "Empty filter", filter: nil, parts: parts, expected: parts},
		{name: "Filter by one category", filter: []model.Category{model.CategoryEngine}, parts: parts, expected: []repoModel.Part{parts[0]}},
		{name: "Filter by two categories", filter: []model.Category{model.CategoryFuel, model.CategoryWing}, parts: parts, expected: []repoModel.Part{parts[1], parts[2]}},
		{name: "Filter by non-existing category", filter: []model.Category{model.CategoryUnspecified}, parts: parts, expected: []repoModel.Part{}},
		{name: "Filter by one exist and one non-existing", filter: []model.Category{model.CategoryWing, model.CategoryUnspecified}, parts: parts, expected: []repoModel.Part{parts[2]}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := filterByCategory(tt.filter, tt.parts)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestFilterByCountry(t *testing.T) {
	parts := []repoModel.Part{
		{Name: "Main engine", Manufacturer: repoModel.Manufacturer{Country: "USA"}},
		{Name: "Fuel tank", Manufacturer: repoModel.Manufacturer{Country: "Germany"}},
		{Name: "Left wing", Manufacturer: repoModel.Manufacturer{Country: "Japan"}},
		{Name: "Tail", Manufacturer: repoModel.Manufacturer{Name: "USA", Country: "Japan"}},
	}

	tests := []struct {
		name     string
		filter   []string
		parts    []repoModel.Part
		expected []repoModel.Part
	}{
		{name: "Empty filter", filter: nil, parts: parts, expected: parts},
		{name: "Filter by one country", filter: []string{"USA"}, parts: parts, expected: []repoModel.Part{parts[0]}},
		{name: "Filter by two countries", filter: []string{"USA", "Japan"}, parts: parts, expected: []repoModel.Part{parts[0], parts[2], parts[3]}},
		{name: "Filter by non-existing countries", filter: []string{"France", "Italy"}, parts: parts, expected: []repoModel.Part{}},
		{name: "Filter by one exist and one non-existing", filter: []string{"USA", "France"}, parts: parts, expected: []repoModel.Part{parts[0]}},
		{name: "One country matches two parts", filter: []string{"Japan"}, parts: parts, expected: []repoModel.Part{parts[2], parts[3]}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := filterByCountry(tt.filter, tt.parts)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestFilterByTag(t *testing.T) {
	parts := []repoModel.Part{
		{Name: "Main engine", Tags: []string{"critical", "propulsion"}},
		{Name: "Fuel tank", Tags: []string{"high-risk", "flammable"}},
		{Name: "Left wing", Tags: []string{"aerodynamic", "left-side"}},
		{Name: "Tail", Tags: []string{"steering", "aerodynamic", "rear"}},
		{Name: "Bolt", Tags: nil},
	}

	tests := []struct {
		name     string
		filter   []string
		parts    []repoModel.Part
		expected []repoModel.Part
	}{
		{name: "Empty filter", filter: nil, parts: parts, expected: parts},
		{name: "Filter by one tag", filter: []string{"critical"}, parts: parts, expected: []repoModel.Part{parts[0]}},
		{name: "No duplicates on multiple tag match", filter: []string{"aerodynamic", "steering"}, parts: parts, expected: []repoModel.Part{parts[2], parts[3]}},
		{name: "Filter by non-existing tag", filter: []string{"nonexistent"}, parts: parts, expected: []repoModel.Part{}},
		{name: "One tag matches two parts", filter: []string{"aerodynamic"}, parts: parts, expected: []repoModel.Part{parts[2], parts[3]}},
		{name: "Part with no tags is skipped", filter: []string{"critical"}, parts: parts, expected: []repoModel.Part{parts[0]}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := filterByTag(tt.filter, tt.parts)
			require.Equal(t, tt.expected, result)
		})
	}
}
