package converter

import (
	"log"
	"time"

	inventoryV1 "github.com/T1mofey4/microservices/shared/pkg/proto/inventory/v1"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/T1mofey4/microservices/inventory/internal/model"
)

func PartToProto(p *model.Part) *inventoryV1.Part {
	return &inventoryV1.Part{
		Uuid:          p.UUID,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		StockQuantity: p.StockQuantity,
		Category:      categoryToProto(p.Category),
		Dimensions:    dimensionsToProto(p.Dimensions),
		Manufacturer:  manufacturerToProto(p.Manufacturer),
		Tags:          p.Tags,
		CreatedAt:     timestamppb.New(p.CreatedAt),
		UpdatedAt:     timeToProto(p.UpdatedAt),
	}
}

func categoryToProto(c model.Category) inventoryV1.Category {
	switch c {
	case model.CategoryEngine:
		return inventoryV1.Category_CATEGORY_ENGINE
	case model.CategoryFuel:
		return inventoryV1.Category_CATEGORY_FUEL
	case model.CategoryPorthole:
		return inventoryV1.Category_CATEGORY_PORTHOLE
	case model.CategoryWing:
		return inventoryV1.Category_CATEGORY_WING
	case model.CategoryUnspecified:
		return inventoryV1.Category_CATEGORY_UNSPECIFIED
	default:
		log.Printf("unknown category: %d", c)
		return inventoryV1.Category_CATEGORY_UNSPECIFIED
	}
}

func dimensionsToProto(d model.Dimensions) *inventoryV1.Dimensions {
	return &inventoryV1.Dimensions{
		Length: d.Length,
		Width:  d.Width,
		Height: d.Height,
		Weight: d.Weight,
	}
}

func manufacturerToProto(m model.Manufacturer) *inventoryV1.Manufacturer {
	return &inventoryV1.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

func timeToProto(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}

	return timestamppb.New(*t)
}

func FilterToModel(f *inventoryV1.PartsFilter) model.PartsFilter {
	if f == nil {
		return model.PartsFilter{}
	}

	return model.PartsFilter{
		UUIDs:                 f.Uuids,
		Names:                 f.Names,
		Categories:            categoriesToModel(f.Categories),
		ManufacturerCountries: f.ManufacturerCountries,
		Tags:                  f.Tags,
	}
}

func categoriesToModel(categories []inventoryV1.Category) []model.Category {
	categoryModel := make([]model.Category, 0, len(categories))
	for _, c := range categories {
		categoryModel = append(categoryModel, categoryFromProto(c))
	}

	return categoryModel
}

func categoryFromProto(c inventoryV1.Category) model.Category {
	switch c {
	case inventoryV1.Category_CATEGORY_ENGINE:
		return model.CategoryEngine
	case inventoryV1.Category_CATEGORY_FUEL:
		return model.CategoryFuel
	case inventoryV1.Category_CATEGORY_PORTHOLE:
		return model.CategoryPorthole
	case inventoryV1.Category_CATEGORY_WING:
		return model.CategoryWing
	case inventoryV1.Category_CATEGORY_UNSPECIFIED:
		return model.CategoryUnspecified
	default:
		log.Printf("categoryFromProto: unknown category %d", c)
		return model.CategoryUnspecified
	}
}

func PartsToProto(parts []model.Part) []*inventoryV1.Part {
	partModel := make([]*inventoryV1.Part, 0, len(parts))
	for _, p := range parts {
		partModel = append(partModel, PartToProto(&p))
	}

	return partModel
}
