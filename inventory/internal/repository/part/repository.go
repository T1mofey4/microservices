package part

import (
	"sync"

	def "github.com/T1mofey4/microservices/inventory/internal/repository"
	repoModel "github.com/T1mofey4/microservices/inventory/internal/repository/model"
)

var _ def.Repository = (*repository)(nil)

type repository struct {
	mu    sync.RWMutex
	parts map[string]repoModel.Part
}

func NewRepository() *repository {
	return &repository{
		parts: make(map[string]repoModel.Part),
	}
}
