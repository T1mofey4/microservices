package part

import (
	"github.com/T1mofey4/microservices/inventory/internal/repository"
	def "github.com/T1mofey4/microservices/inventory/internal/service"
)

var _ def.PartService = (*partService)(nil)

type partService struct {
	repo repository.PartRepository
}

func NewService(repo repository.PartRepository) *partService {
	return &partService{
		repo: repo,
	}
}
