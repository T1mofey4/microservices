package v1

import (
	"context"
	"fmt"

	inventoryV1 "github.com/T1mofey4/microservices/shared/pkg/proto/inventory/v1"
	"github.com/google/uuid"

	"github.com/T1mofey4/microservices/order/internal/client/converter"
	"github.com/T1mofey4/microservices/order/internal/model"
)

// ListParts получает из inventory список запчастей
func (c *inventoryClient) ListParts(ctx context.Context, partUUIDs []uuid.UUID) ([]model.PartInfo, error) {
	partUUID := make([]string, 0, len(partUUIDs))
	for _, u := range partUUIDs {
		partUUID = append(partUUID, u.String())
	}

	filter := inventoryV1.PartsFilter{Uuids: partUUID}
	req := inventoryV1.ListPartsRequest{Filter: &filter}
	res, err := c.client.ListParts(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("inventory ListParts: %w", err)
	}

	if len(res.Parts) != len(partUUID) {
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
