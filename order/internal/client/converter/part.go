package converter

import (
	"fmt"

	inventoryV1 "github.com/T1mofey4/microservices/shared/pkg/proto/inventory/v1"
	"github.com/google/uuid"

	"github.com/T1mofey4/microservices/order/internal/model"
)

func PartToModel(part *inventoryV1.Part) (model.PartInfo, error) {
	partUUID, err := uuid.Parse(part.Uuid)
	if err != nil {
		return model.PartInfo{}, fmt.Errorf("incorrect part uuid: %w", err)
	}

	return model.PartInfo{
		PartUUID: partUUID,
		Price:    part.Price,
	}, nil
}
