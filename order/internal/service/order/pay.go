package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/T1mofey4/microservices/order/internal/model"
	repo "github.com/T1mofey4/microservices/order/internal/repository"
)

func (s *orderService) Pay(ctx context.Context, orderUUID uuid.UUID, method model.PaymentMethod) (uuid.UUID, error) {
	order, err := s.orderRepo.Get(ctx, orderUUID.String())
	if err != nil {
		return uuid.Nil, fmt.Errorf("get order failed: %w", err)
	}

	if order.Status == model.StatusPaid {
		return uuid.Nil, fmt.Errorf("status already paid: %w", model.ErrOrderAlreadyPaid)
	}

	if order.Status == model.StatusCancelled {
		return uuid.Nil, fmt.Errorf("status already cancelled: %w", model.ErrOrderCancelled)
	}

	if order.Status != model.StatusPendingPayment {
		return uuid.Nil, model.ErrInvalidArgument
	}

	transactionUUID, err := s.paymentClient.PayOrder(ctx, order.UserUUID, orderUUID, method)
	if err != nil {
		return uuid.Nil, fmt.Errorf("payment failed: %w", err)
	}

	updOrder := model.Order{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUUIDs:       order.PartUUIDs,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: &transactionUUID,
		PaymentMethod:   &method,
		Status:          model.StatusPaid,
	}

	if err = s.orderRepo.Update(ctx, orderUUID.String(), model.StatusPendingPayment, updOrder); err != nil {
		switch {
		case errors.Is(err, model.ErrOrderNotFound):
			return uuid.Nil, err
		case errors.Is(err, repo.ErrStatusConflict):
			//получаем заказ по uuid еще раз и проверяем какой именно статус у заказа
			retryOrder, retryErr := s.orderRepo.Get(ctx, orderUUID.String())
			if retryErr != nil {
				return uuid.Nil, fmt.Errorf("retry get order failed: %w", retryErr)
			}

			switch retryOrder.Status {
			case model.StatusPaid:
				return uuid.Nil, fmt.Errorf("order already paid: %w", model.ErrOrderAlreadyPaid)
			case model.StatusCancelled:
				return uuid.Nil, fmt.Errorf("order cancelled: %w", model.ErrOrderCancelled)
			default:
				return uuid.Nil, fmt.Errorf("failed to update order: %w", repo.ErrStatusConflict)
			}

		default:
			return uuid.Nil, fmt.Errorf("failed to update order status: %w", err)
		}
	}

	return transactionUUID, nil
}
