package service

import (
	"context"

	"github.com/T1mofey4/microservices/inventory/internal/model"
)

type Service interface {
	GetPart(ctx context.Context, uuid string) (*model.Part, error)
	ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error)
}
