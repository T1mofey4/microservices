package part

import (
	"context"

	"github.com/T1mofey4/microservices/inventory/internal/model"
	"github.com/T1mofey4/microservices/inventory/internal/repository/converter"
)

func (r *partRepository) Get(ctx context.Context, uuid string) (*model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.parts[uuid]
	if !ok {
		return nil, model.ErrPartNotFound
	}

	return converter.PartToModel(p), nil
}
