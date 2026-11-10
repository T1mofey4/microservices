package order

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"

	def "github.com/T1mofey4/microservices/order/internal/repository"
)

var _ def.OrderRepository = (*orderRepository)(nil)

type orderRepository struct {
	pool    *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewOrderRepository(pool *pgxpool.Pool) *orderRepository {
	return &orderRepository{
		pool:    pool,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}
