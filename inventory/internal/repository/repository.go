package repository

import (
	"context"

	"github.com/T1mofey4/microservices/inventory/internal/model"
)

type PartRepository interface {
	Get(ctx context.Context, uuid string) (*model.Part, error)
	ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error)
}
