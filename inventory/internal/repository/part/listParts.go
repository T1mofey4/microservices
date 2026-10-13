package part

import (
	"context"

	"github.com/T1mofey4/microservices/inventory/internal/model"
)

func (r *repository) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	//	parts := make

	return []model.Part{}, nil

}
