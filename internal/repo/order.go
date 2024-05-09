package repo

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sprint-id/eniqilo-server/internal/dto"
	"github.com/sprint-id/eniqilo-server/internal/entity"
	"github.com/sprint-id/eniqilo-server/internal/ierr"
	timepkg "github.com/sprint-id/eniqilo-server/pkg/time"
)

type orderRepo struct {
	conn *pgxpool.Pool
}

func newOrderRepo(conn *pgxpool.Pool) *orderRepo {
	return &orderRepo{conn}
}

func (mr *orderRepo) AddOrder(ctx context.Context, sub string, order entity.Order) (dto.ResOrder, error) {
	// get price and stock
	// if not found, return error
	var total int
	for _, pd := range order.ProductDetails {
		q := `SELECT id, price, stock FROM products WHERE id = $1`
		var id string
		var price int
		var stock int
		err := mr.conn.QueryRow(ctx, q, pd.ProductID).Scan(&id, &price, &stock)
		if err != nil {
			return dto.ResOrder{}, ierr.ErrNotFound
		}

		// check stock
		if stock < pd.Quantity {
			return dto.ResOrder{}, ierr.ErrStockNotEnough
		}

		total += price * pd.Quantity
	}
	// check paid, if not enough return error
	if total > order.Paid {
		return dto.ResOrder{}, ierr.ErrNotEnoughPaid
	}

	// check change
	if order.Change < 0 {
		return dto.ResOrder{}, ierr.ErrBadRequest
	}

	// insert order
	q := `INSERT INTO orders (customer_id, paid, change, created_at)
	VALUES ( $1, $2, $3, $4, $5, EXTRACT(EPOCH FROM now())::bigint) RETURNING id, created_at`

	var orderId string
	var createdAt time.Time
	err := mr.conn.QueryRow(ctx, q, sub, order.CustomerID, order.Paid, order.Change).Scan(&orderId, &createdAt)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return dto.ResOrder{}, ierr.ErrDuplicate
			}
		}
		return dto.ResOrder{}, err
	}

	// insert all order details
	for _, pd := range order.ProductDetails {
		q := `INSERT INTO order_details (order_id, product_id, quantity, created_at)
		VALUES ( $1, $2, $3, EXTRACT(EPOCH FROM now())::bigint) RETURNING id`

		var id string
		err := mr.conn.QueryRow(ctx, q, orderId, pd.ProductID, pd.Quantity).Scan(&id)
		if err != nil {
			if pgErr, ok := err.(*pgconn.PgError); ok {
				if pgErr.Code == "23505" {
					return dto.ResOrder{}, ierr.ErrDuplicate
				}
			}
			return dto.ResOrder{}, err
		}
	}

	return dto.ResOrder{ID: orderId, CreatedAt: timepkg.TimeToISO8601(createdAt)}, nil
}
