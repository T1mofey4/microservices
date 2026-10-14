package part

import (
	"context"

	"github.com/T1mofey4/microservices/inventory/internal/model"
)

func (s *service) GetPart(ctx context.Context, uuid string) (*model.Part, error) {
	part, err := s.repo.Get(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return part, nil
}
