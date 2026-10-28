package part

import (
	"sync"

	def "github.com/T1mofey4/microservices/inventory/internal/repository"
	repoModel "github.com/T1mofey4/microservices/inventory/internal/repository/model"
)

var _ def.PartRepository = (*partRepository)(nil)

type partRepository struct {
	mu    sync.RWMutex
	parts map[string]repoModel.Part
}

func NewPartRepository() *partRepository {
	tempParts := initParts()
	parts := make(map[string]repoModel.Part, len(tempParts))
	for _, p := range tempParts {
		parts[p.UUID] = p
	}

	return &partRepository{
		parts: parts,
	}
}
