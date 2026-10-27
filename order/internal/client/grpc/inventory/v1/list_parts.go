package v1

import (
	"context"
	"fmt"

	inventoryV1 "github.com/T1mofey4/microservices/shared/pkg/proto/inventory/v1"

	"github.com/T1mofey4/microservices/order/internal/client/converter"
	"github.com/T1mofey4/microservices/order/internal/model"
)

// ListParts получает из inventory список запчастей
func (c *inventoryClient) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.PartInfo, error) {
	req := inventoryV1.ListPartsRequest{Filter: converter.FilterToProto(filter)}
	res, err := c.client.ListParts(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("inventory ListParts: %w", err)
	}

	if len(res.Parts) != len(filter.UUIDs) {
		return nil, model.ErrPartNotFound
	}

	result := make([]model.PartInfo, 0, len(res.Parts))
	for _, p := range res.Parts {
		part, err := converter.PartToModel(p)
		if err != nil {
			return nil, fmt.Errorf("inventory ListParts: %w", err)
		}
		result = append(result, part)
	}

	return result, nil
}
