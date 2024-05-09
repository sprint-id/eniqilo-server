package repo

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	conn *pgxpool.Pool

	User        *userRepo
	Product     *productRepo
	Customer    *customerRepo
	Transaction *transactionRepo
}

func NewRepo(conn *pgxpool.Pool) *Repo {
	repo := Repo{}
	repo.conn = conn

	repo.User = newUserRepo(conn)
	repo.Product = newProductRepo(conn)
	repo.Customer = newCustomerRepo(conn)
	repo.Transaction = newOrderRepo(conn)

	return &repo
}
