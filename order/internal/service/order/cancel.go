package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/T1mofey4/microservices/order/internal/model"
	repo "github.com/T1mofey4/microservices/order/internal/repository"
)

func (s *orderService) Cancel(ctx context.Context, orderUUID uuid.UUID) error {
	order, err := s.orderRepo.Get(ctx, orderUUID.String())
	if err != nil {
		return fmt.Errorf("get order failed: %w", err)
	}

	if order.Status == model.StatusPaid {
		return model.ErrOrderAlreadyPaid
	}

	if order.Status == model.StatusCancelled {
		return model.ErrOrderCancelled
	}

	if order.Status != model.StatusPendingPayment {
		return model.ErrInvalidArgument
	}

	updOrder := order
	updOrder.Status = model.StatusCancelled

	if err := s.orderRepo.Update(ctx, orderUUID.String(), model.StatusPendingPayment, updOrder); err != nil {
		switch {
		case errors.Is(err, model.ErrOrderNotFound):
			return err
		case errors.Is(err, repo.ErrStatusConflict):
			//получаем заказ по uuid еще раз и проверяем какой именно статус у заказа
			retryOrder, retryErr := s.orderRepo.Get(ctx, orderUUID.String())
			if retryErr != nil {
				return fmt.Errorf("retry get order failed: %w", retryErr)
			}

			switch retryOrder.Status {
			case model.StatusPaid:
				return fmt.Errorf("order already paid: %w", model.ErrOrderAlreadyPaid)
			case model.StatusCancelled:
				return fmt.Errorf("order already cancelled: %w", model.ErrOrderCancelled)
			default:
				return fmt.Errorf("failed to update order: %w", repo.ErrStatusConflict)
			}

		default:
			return fmt.Errorf("failed to update order status: %w", err)
		}
	}

	return nil
}
