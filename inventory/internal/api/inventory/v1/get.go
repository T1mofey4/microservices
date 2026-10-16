package v1

import (
	"context"
	"errors"
	"log"

	"github.com/T1mofey4/microservices/inventory/internal/converter"
	"github.com/T1mofey4/microservices/inventory/internal/model"
	inventoryV1 "github.com/T1mofey4/microservices/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	part, err := a.service.GetPart(ctx, req.Uuid)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrPartNotFound):
			return nil, status.Errorf(codes.NotFound, "деталь uuid: %s не найдена", req.Uuid)
		default:
			log.Printf("GetPart failed: %v UUID: %s", err, req.Uuid)
			return nil, status.Error(codes.Internal, "внутренняя ошибка")
		}
	}

	partToProto := converter.PartToProto(part)

	return &inventoryV1.GetPartResponse{Part: partToProto}, nil
}
