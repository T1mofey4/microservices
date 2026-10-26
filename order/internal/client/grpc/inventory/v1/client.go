package v1

import (
	inventoryV1 "github.com/T1mofey4/microservices/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"

	def "github.com/T1mofey4/microservices/order/internal/client/grpc"
)

var _ def.InventoryClient = (*inventoryClient)(nil)

type inventoryClient struct {
	client inventoryV1.InventoryServiceClient
}

func NewInventoryClient(conn *grpc.ClientConn) *inventoryClient {
	return &inventoryClient{
		client: inventoryV1.NewInventoryServiceClient(conn),
	}
}
