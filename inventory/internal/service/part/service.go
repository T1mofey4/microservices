package part

import (
	"github.com/T1mofey4/microservices/inventory/internal/repository"
	def "github.com/T1mofey4/microservices/inventory/internal/service"
)

var _ def.Service = (*service)(nil)

type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *service {
	return &service{
		repo: repo,
	}
}
