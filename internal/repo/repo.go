package repo

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	conn *pgxpool.Pool

	User    *userRepo
	Product *productRepo
	Match   *matchRepo
}

func NewRepo(conn *pgxpool.Pool) *Repo {
	repo := Repo{}
	repo.conn = conn

	repo.User = newUserRepo(conn)
	repo.Product = newProductRepo(conn)
	repo.Match = newMatchRepo(conn)

	return &repo
}
