package order

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"github.com/T1mofey4/microservices/order/internal/model"
	"github.com/T1mofey4/microservices/order/internal/repository/converter"
	repoModel "github.com/T1mofey4/microservices/order/internal/repository/model"
)

// Get получает заказ по UUID
func (r *orderRepository) Get(ctx context.Context, orderUUID string) (model.Order, error) {
	builderSelect := r.builder.Select(
		"order_uuid",
		"user_uuid",
		"part_uuids",
		"total_price",
		"transaction_uuid",
		"payment_method",
		"status",
	).
		From("orders").
		Where(sq.Eq{"order_uuid": orderUUID})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return model.Order{}, fmt.Errorf("build select order query: %w", err)
	}

	var repoOrder repoModel.Order
	err = r.pool.QueryRow(ctx, query, args...).Scan(
		&repoOrder.OrderUUID,
		&repoOrder.UserUUID,
		&repoOrder.PartUUIDs,
		&repoOrder.TotalPrice,
		&repoOrder.TransactionUUID,
		&repoOrder.PaymentMethod,
		&repoOrder.Status,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Order{}, fmt.Errorf("get order: %w", model.ErrOrderNotFound)
		}
		return model.Order{}, fmt.Errorf("scan order: %w", err)
	}

	modelOrder, err := converter.OrderToModel(repoOrder)
	if err != nil {
		return model.Order{}, fmt.Errorf("get order %s: %w", orderUUID, err)
	}

	return modelOrder, nil
}
