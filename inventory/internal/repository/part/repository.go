package part

import (
	"go.mongodb.org/mongo-driver/mongo"

	def "github.com/T1mofey4/microservices/inventory/internal/repository"
)

var _ def.PartRepository = (*partRepository)(nil)

type partRepository struct {
	collection *mongo.Collection
}

func NewPartRepository(db *mongo.Database) *partRepository {
	collection := db.Collection("inventory-service")

	return &partRepository{
		collection: collection,
	}
}
