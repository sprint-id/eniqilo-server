package repo

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sprint-id/eniqilo-server/internal/dto"
	"github.com/sprint-id/eniqilo-server/internal/entity"
	"github.com/sprint-id/eniqilo-server/internal/ierr"
)

type customerRepo struct {
	conn *pgxpool.Pool
}

func newCustomerRepo(conn *pgxpool.Pool) *customerRepo {
	return &customerRepo{conn}
}

func (mr *customerRepo) RegisterCustomer(ctx context.Context, sub string, customer entity.Customer) (dto.ResRegisterOrGetCustomer, error) {
	// Start a transaction with serializable isolation level
	tx, err := mr.conn.Begin(ctx)
	if err != nil {
		return dto.ResRegisterOrGetCustomer{}, err
	}
	defer tx.Rollback(ctx)

	q := `INSERT INTO customers (staff_id, phone_number, name, created_at)
	VALUES ( $1, $2, $3, EXTRACT(EPOCH FROM now())::bigint) RETURNING id`

	var id string
	err = tx.QueryRow(ctx, q, sub, customer.PhoneNumber, customer.Name).Scan(&id)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return dto.ResRegisterOrGetCustomer{}, ierr.ErrDuplicate
			}
		}
		return dto.ResRegisterOrGetCustomer{}, err
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		return dto.ResRegisterOrGetCustomer{}, err
	}

	return dto.ResRegisterOrGetCustomer{UserID: id, PhoneNumber: customer.PhoneNumber, Name: customer.Name}, nil
}

// func (mr *customerRepo) RegisterCustomer(ctx context.Context, sub string, customer entity.Customer) (dto.ResRegisterOrGetCustomer, error) {
// 	q := `INSERT INTO customers (staff_id, phone_number, name, created_at)
// 	VALUES ( $1, $2, $3, EXTRACT(EPOCH FROM now())::bigint) RETURNING id`

// 	var id string
// 	err := mr.conn.QueryRow(ctx, q, sub, customer.PhoneNumber, customer.Name).Scan(&id)
// 	if err != nil {
// 		if pgErr, ok := err.(*pgconn.PgError); ok {
// 			if pgErr.Code == "23505" {
// 				return dto.ResRegisterOrGetCustomer{}, ierr.ErrDuplicate
// 			}
// 		}
// 		return dto.ResRegisterOrGetCustomer{}, err
// 	}

// 	return dto.ResRegisterOrGetCustomer{UserID: id, PhoneNumber: customer.PhoneNumber, Name: customer.Name}, nil
// }

func (mr *customerRepo) GetCustomer(ctx context.Context, param dto.ParamGetCustomer, sub string) ([]dto.ResRegisterOrGetCustomer, error) {
	var query strings.Builder

	query.WriteString("SELECT id, phone_number, name FROM customers WHERE 1=1 ")

	// param phone number: it should search by wildcard (ex: if search by phoneNumber=+62 then customer with phone number +628123... should appear, but phoneNumber=123 will not show that)
	if param.PhoneNumber != "" {
		query.WriteString(fmt.Sprintf("AND phone_number LIKE '%s' ", fmt.Sprintf("%%%s%%", param.PhoneNumber)))
	}

	if param.Name != "" {
		query.WriteString(fmt.Sprintf("AND LOWER(name) LIKE LOWER('%s') ", fmt.Sprintf("%%%s%%", param.Name)))
	}

	rows, err := mr.conn.Query(ctx, query.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []dto.ResRegisterOrGetCustomer
	for rows.Next() {
		var customer dto.ResRegisterOrGetCustomer
		err = rows.Scan(&customer.UserID, &customer.PhoneNumber, &customer.Name)
		if err != nil {
			return nil, err
		}

		customers = append(customers, customer)
	}

	return customers, nil
}
