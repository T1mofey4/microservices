package converter

import (
	"log"

	"github.com/T1mofey4/microservices/inventory/internal/model"
	repoModel "github.com/T1mofey4/microservices/inventory/internal/repository/model"
)

func PartToModel(p repoModel.Part) *model.Part {
	return &model.Part{
		UUID:          p.UUID,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		StockQuantity: p.StockQuantity,
		Category:      categoryToModel(p.Category),
		Dimensions:    dimensionsToModel(p.Dimensions),
		Manufacturer:  manufacturerToModel(p.Manufacturer),
		Tags:          p.Tags,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}

func categoryToModel(c repoModel.Category) model.Category {
	switch c {
	case repoModel.CategoryEngine:
		return model.CategoryEngine
	case repoModel.CategoryFuel:
		return model.CategoryFuel
	case repoModel.CategoryPorthole:
		return model.CategoryPorthole
	case repoModel.CategoryWing:
		return model.CategoryWing
	case repoModel.CategoryUnspecified:
		return model.CategoryUnspecified
	default:
		log.Printf("unknown category: %d", c)
		return model.CategoryUnspecified
	}
}

func CategoryToRepoModel(c model.Category) repoModel.Category {
	switch c {
	case model.CategoryEngine:
		return repoModel.CategoryEngine
	case model.CategoryFuel:
		return repoModel.CategoryFuel
	case model.CategoryPorthole:
		return repoModel.CategoryPorthole
	case model.CategoryWing:
		return repoModel.CategoryWing
	case model.CategoryUnspecified:
		return repoModel.CategoryUnspecified
	default:
		log.Printf("unknown category: %d", c)
		return repoModel.CategoryUnspecified
	}
}

func dimensionsToModel(d repoModel.Dimensions) model.Dimensions {
	return model.Dimensions{
		Length: d.Length,
		Width:  d.Width,
		Height: d.Height,
		Weight: d.Weight,
	}
}

func manufacturerToModel(m repoModel.Manufacturer) model.Manufacturer {
	return model.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

func PartsToModel(parts []repoModel.Part) []model.Part {
	result := make([]model.Part, 0, len(parts))
	for _, p := range parts {
		result = append(result, *PartToModel(p))
	}
	return result
}
