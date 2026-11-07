package v1

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/T1mofey4/microservices/inventory/internal/converter"
	inventoryV1 "github.com/T1mofey4/microservices/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	parts, err := a.service.ListParts(ctx, converter.FilterToModel(req.Filter))
	if err != nil {
		log.Printf("ListParts: failed %v", err)
		return nil, status.Error(codes.Internal, "внутренняя ошибка")
	}

	partsToProto := converter.PartsToProto(parts)

	return &inventoryV1.ListPartsResponse{Parts: partsToProto}, nil
}
