package part

import (
	"context"

	"github.com/T1mofey4/microservices/inventory/internal/model"
)

func (r *repository) Get(ctx context.Context, uuid string) (*model.Part, error) {
	return &model.Part{}, nil
}
