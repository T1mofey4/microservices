package part

import (
	"context"

	"github.com/T1mofey4/microservices/inventory/internal/model"
)

func (s *service) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	parts, err := s.repo.ListParts(ctx, filter)
	if err != nil {
		return nil, err
	}
	return parts, nil
}
