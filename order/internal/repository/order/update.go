package order

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/T1mofey4/microservices/order/internal/model"
	def "github.com/T1mofey4/microservices/order/internal/repository"
	"github.com/T1mofey4/microservices/order/internal/repository/converter"
)

// Update обновляет заказ
func (r *orderRepository) Update(ctx context.Context, orderUUID string, expectedStatus model.OrderStatus, updOrder model.Order) error {
	updRepoOrder := converter.OrderToRepoModel(updOrder)

	builderUpdate := r.builder.
		Update("orders").
		Set("transaction_uuid", updRepoOrder.TransactionUUID).
		Set("payment_method", updRepoOrder.PaymentMethod).
		Set("status", updRepoOrder.Status).
		Where(sq.Eq{
			"order_uuid": orderUUID,
			"status":     converter.StatusToRepoModel(expectedStatus),
		})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		return fmt.Errorf("build update order query: %w", err)
	}

	res, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("exec update order: %w", err)
	}

	rows := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("update order %s: %w", orderUUID, def.ErrStatusConflict)
	}

	return nil
}
