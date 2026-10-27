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

func FilterToProto(filter model.PartsFilter) *inventoryV1.PartsFilter {
	uuids := make([]string, 0, len(filter.UUIDs))
	for _, partUUID := range filter.UUIDs {
		uuids = append(uuids, partUUID.String())
	}

	categories := make([]inventoryV1.Category, 0, len(filter.Categories))
	for _, c := range filter.Categories {
		switch c {
		case model.CategoryUnspecified:
			categories = append(categories, inventoryV1.Category_CATEGORY_UNSPECIFIED)
		case model.CategoryEngine:
			categories = append(categories, inventoryV1.Category_CATEGORY_ENGINE)
		case model.CategoryFuel:
			categories = append(categories, inventoryV1.Category_CATEGORY_FUEL)
		case model.CategoryPorthole:
			categories = append(categories, inventoryV1.Category_CATEGORY_PORTHOLE)
		case model.CategoryWing:
			categories = append(categories, inventoryV1.Category_CATEGORY_WING)
		default:
			categories = append(categories, inventoryV1.Category_CATEGORY_UNSPECIFIED)
		}

	}

	return &inventoryV1.PartsFilter{
		Uuids:                 uuids,
		Names:                 filter.Names,
		Categories:            categories,
		Tags:                  filter.Tags,
		ManufacturerCountries: filter.ManufacturerCountries,
	}
}
