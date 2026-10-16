package v1

import (
	inventoryV1 "github.com/T1mofey4/microservices/shared/pkg/proto/inventory/v1"

	"github.com/T1mofey4/microservices/inventory/internal/service"
)

var _ inventoryV1.InventoryServiceServer = (*api)(nil)

type api struct {
	inventoryV1.UnimplementedInventoryServiceServer
	service service.Service
}

func NewApi(inventoryService service.Service) *api {
	return &api{
		service: inventoryService,
	}
}
