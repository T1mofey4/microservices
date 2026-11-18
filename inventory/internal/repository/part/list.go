package part

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/T1mofey4/microservices/inventory/internal/model"
	"github.com/T1mofey4/microservices/inventory/internal/repository/converter"
	repoModel "github.com/T1mofey4/microservices/inventory/internal/repository/model"
)

func (r *partRepository) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	mongoFilter := buildListFilter(filter)

	cursor, err := r.collection.Find(ctx, mongoFilter)
	if err != nil {
		return nil, fmt.Errorf("find parts: %w", err)
	}
	defer func() {
		cerr := cursor.Close(ctx)
		if cerr != nil {
			log.Printf("failed to close cursor: %v\n", cerr)
		}
	}()

	var parts []repoModel.Part
	err = cursor.All(ctx, &parts)
	if err != nil {
		return nil, fmt.Errorf("decode parts: %w", err)
	}

	return converter.PartsToModel(parts), nil
}

func buildListFilter(filter model.PartsFilter) bson.M {
	list := bson.M{}

	if len(filter.UUIDs) > 0 {
		list["uuid"] = bson.M{"$in": filter.UUIDs}
	}

	if len(filter.Names) > 0 {
		list["name"] = bson.M{"$in": filter.Names}
	}

	if len(filter.Categories) > 0 {
		repoCategories := make([]repoModel.Category, 0, len(filter.Categories))
		for _, c := range filter.Categories {
			repoCategories = append(repoCategories, converter.CategoryToRepoModel(c))
		}

		list["category"] = bson.M{"$in": repoCategories}
	}

	if len(filter.ManufacturerCountries) > 0 {
		list["manufacturer.country"] = bson.M{"$in": filter.ManufacturerCountries}
	}

	if len(filter.Tags) > 0 {
		list["tags"] = bson.M{"$in": filter.Tags}
	}

	return list
}
