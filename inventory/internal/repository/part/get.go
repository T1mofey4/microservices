package part

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/T1mofey4/microservices/inventory/internal/model"
	"github.com/T1mofey4/microservices/inventory/internal/repository/converter"
	repoModel "github.com/T1mofey4/microservices/inventory/internal/repository/model"
)

func (r *partRepository) Get(ctx context.Context, uuid string) (*model.Part, error) {
	var part repoModel.Part
	if err := r.collection.FindOne(ctx, bson.M{"uuid": uuid}).Decode(&part); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, model.ErrPartNotFound
		}

		return nil, fmt.Errorf("get part: %w", err)
	}

	return converter.PartToModel(part), nil
}
