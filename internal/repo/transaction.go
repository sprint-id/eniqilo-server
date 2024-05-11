package repo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sprint-id/eniqilo-server/internal/dto"
	"github.com/sprint-id/eniqilo-server/internal/entity"
	"github.com/sprint-id/eniqilo-server/internal/ierr"
	timepkg "github.com/sprint-id/eniqilo-server/pkg/time"
)

type transactionRepo struct {
	conn *pgxpool.Pool
}

func newTransactionRepo(conn *pgxpool.Pool) *transactionRepo {
	return &transactionRepo{conn}
}

func (tr *transactionRepo) AddTransaction(ctx context.Context, sub string, transaction entity.Transaction) (dto.ResTransaction, error) {
	// get price and stock
	// if not found, return error
	var total int
	for _, pd := range transaction.ProductDetails {
		q := `SELECT id, price, stock FROM products WHERE id = $1`
		var id string
		var price int
		var stock int
		err := tr.conn.QueryRow(ctx, q, pd.ProductID).Scan(&id, &price, &stock)
		if err != nil {
			return dto.ResTransaction{}, ierr.ErrNotFound
		}

		// check stock
		if stock < pd.Quantity {
			return dto.ResTransaction{}, ierr.ErrStockNotEnough
		}

		total += price * pd.Quantity
	}
	// check paid, if not enough return error
	if total > transaction.Paid {
		return dto.ResTransaction{}, ierr.ErrNotEnoughPaid
	}

	// check change
	if transaction.Change < 0 {
		return dto.ResTransaction{}, ierr.ErrBadRequest
	} else if transaction.Change != total-transaction.Paid {
		return dto.ResTransaction{}, ierr.ErrChangeNotMatch
	}

	// insert transaction
	q := `INSERT INTO transactions (customer_id, paid, change, created_at)
	VALUES ( $1, $2, $3, EXTRACT(EPOCH FROM now())::bigint) RETURNING id, created_at`

	var transactionId string
	var createdAt int64
	err := tr.conn.QueryRow(ctx, q, transaction.CustomerID, transaction.Paid, transaction.Change).Scan(&transactionId, &createdAt)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return dto.ResTransaction{}, ierr.ErrDuplicate
			}
		}
		return dto.ResTransaction{}, err
	}

	// insert all transaction details
	for _, pd := range transaction.ProductDetails {
		q := `INSERT INTO transaction_details (transaction_id, product_id, quantity, created_at)
		VALUES ( $1, $2, $3, EXTRACT(EPOCH FROM now())::bigint) RETURNING id`

		var id string
		err := tr.conn.QueryRow(ctx, q, transactionId, pd.ProductID, pd.Quantity).Scan(&id)
		if err != nil {
			if pgErr, ok := err.(*pgconn.PgError); ok {
				if pgErr.Code == "23505" {
					return dto.ResTransaction{}, ierr.ErrDuplicate
				}
			}
			return dto.ResTransaction{}, err
		}
	}

	// Update quantity of product
	for _, pd := range transaction.ProductDetails {
		q := `UPDATE products SET stock = stock - $1 WHERE id = $2`
		_, err := tr.conn.Exec(ctx, q, pd.Quantity, pd.ProductID)
		if err != nil {
			return dto.ResTransaction{}, err
		}
	}

	createdAtTime := time.Unix(createdAt, 0)
	return dto.ResTransaction{ID: transactionId, CreatedAt: timepkg.TimeToISO8601(createdAtTime)}, nil
}

func (tr *transactionRepo) GetTransactionHistory(ctx context.Context, param dto.ParamGetTransactionHistory, sub string) ([]dto.ResTransactionHistory, error) {
	var query strings.Builder
	query.WriteString("SELECT t.id, t.customer_id, t.paid, t.change, t.created_at FROM transactions t WHERE 1=1 ")

	if param.CustomerID != "" {
		query.WriteString(fmt.Sprintf("AND t.customer_id = %s ", param.CustomerID))
	}

	// createdAt sort by created time enum of ASC and DESC
	if param.CreatedAt == "asc" {
		query.WriteString("ORDER BY created_at ASC ")
	} else if param.CreatedAt == "desc" {
		query.WriteString("ORDER BY created_at DESC ")
	}

	// limit and offset
	if param.Limit == 0 {
		param.Limit = 5
	}

	query.WriteString(fmt.Sprintf("LIMIT %d OFFSET %d", param.Limit, param.Offset))

	rows, err := tr.conn.Query(ctx, query.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []dto.ResTransactionHistory
	for rows.Next() {
		var r dto.ResTransactionHistory
		var createdAt int64
		err := rows.Scan(&r.TransactionID, &r.CustomerID, &r.Paid, &r.Change, &createdAt)
		if err != nil {
			return nil, err
		}

		// get product details
		q := `SELECT td.product_id, td.quantity
		FROM transaction_details td
		WHERE td.transaction_id = $1`

		rows, err := tr.conn.Query(ctx, q, r.TransactionID)
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			var pd dto.ReqProductDetail
			err := rows.Scan(&pd.ProductID, &pd.Quantity)
			if err != nil {
				return nil, err
			}
			r.ProductDetails = append(r.ProductDetails, pd)
		}

		createdAtTime := time.Unix(createdAt, 0)
		r.CreatedAt = timepkg.TimeToISO8601(createdAtTime)
		res = append(res, r)
	}

	return res, nil
}
