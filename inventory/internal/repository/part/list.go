package part

import (
	"context"

	"github.com/T1mofey4/microservices/inventory/internal/model"
	"github.com/T1mofey4/microservices/inventory/internal/repository/converter"
	repoModel "github.com/T1mofey4/microservices/inventory/internal/repository/model"
)

func (r *repository) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	parts := make([]repoModel.Part, 0, len(r.parts))
	for _, p := range r.parts {
		parts = append(parts, p)
	}

	result := filterByUUID(filter.UUIDs, parts)
	result = filterByName(filter.Names, result)
	result = filterByCategory(filter.Categories, result)
	result = filterByCountry(filter.ManufacturerCountries, result)
	result = filterByTag(filter.Tags, result)

	return converter.PartsToModel(result), nil
}
