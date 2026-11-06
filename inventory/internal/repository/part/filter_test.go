package part

import (
	"testing"

	"github.com/stretchr/testify/require"

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
