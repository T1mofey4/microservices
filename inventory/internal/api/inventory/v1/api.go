package v1

import (
	"github.com/T1mofey4/microservices/inventory/internal/service"
	inventoryV1 "github.com/T1mofey4/microservices/shared/pkg/proto/inventory/v1"
)

var _ inventoryV1.InventoryServiceServer = (*api)(nil)

type api struct {
	inventoryV1.UnimplementedInventoryServiceServer
	service service.PartService
}

func NewApi(inventoryService service.PartService) *api {
	return &api{
		service: inventoryService,
	}
}
