package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/T1mofey4/microservices/order/internal/model"
	def "github.com/T1mofey4/microservices/order/internal/repository"
	"github.com/T1mofey4/microservices/order/internal/repository/converter"
)

// Create создает заказ
func (r *orderRepository) Create(ctx context.Context, order model.Order) error {
	repoOrder := converter.OrderToRepoModel(order)

	builderInsert := r.builder.
		Insert("orders").
		Columns(
			"order_uuid",
			"user_uuid",
			"part_uuids",
			"total_price",
			"transaction_uuid",
			"payment_method",
			"status",
		).
		Values(
			repoOrder.OrderUUID,
			repoOrder.UserUUID,
			repoOrder.PartUUIDs,
			repoOrder.TotalPrice,
			repoOrder.TransactionUUID,
			repoOrder.PaymentMethod,
			repoOrder.Status,
		)

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return fmt.Errorf("build insert order query: %w", err)
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("create order: %w", def.ErrOrderAlreadyExists)
		}
		return fmt.Errorf("exec insert order: %w", err)
	}

	return nil
}
