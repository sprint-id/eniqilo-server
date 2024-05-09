package repo

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sprint-id/eniqilo-server/internal/entity"
	"github.com/sprint-id/eniqilo-server/internal/ierr"
)

type userRepo struct {
	conn *pgxpool.Pool
}

func newUserRepo(conn *pgxpool.Pool) *userRepo {
	return &userRepo{conn}
}

func (u *userRepo) Insert(ctx context.Context, user entity.User) (string, error) {
	credVal := user.PhoneNumber
	q := `INSERT INTO users (id, name, phone_number, password, created_at)
	VALUES (gen_random_uuid(), $1, $2, $3, EXTRACT(EPOCH FROM now())::bigint) RETURNING id`

	var userID string
	err := u.conn.QueryRow(ctx, q,
		user.Name, credVal, user.Password).Scan(&userID)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return "", ierr.ErrDuplicate
			}
		}
		return "", err
	}

	return userID, nil
}

func (u *userRepo) GetByPhone(ctx context.Context, cred string) (entity.User, error) {
	user := entity.User{}
	q := `SELECT id, name, phone_number, password FROM users
	WHERE phone_number = $1`

	var phone sql.NullString

	err := u.conn.QueryRow(ctx,
		q, cred).Scan(&user.ID, &user.Name, &phone, &user.Password)

	user.PhoneNumber = phone.String

	if err != nil {
		if err.Error() == "no rows in result set" {
			return user, ierr.ErrNotFound
		}
		return user, err
	}

	return user, nil
}

func (u *userRepo) GetByID(ctx context.Context, id string) (entity.User, error) {
	user := entity.User{}
	q := `SELECT phone_number, name, password FROM users
	WHERE id = $1`

	var phone sql.NullString

	err := u.conn.QueryRow(ctx,
		q, id).Scan(&phone, &user.Name, &user.Password)

	user.PhoneNumber = phone.String

	if err != nil {
		if err.Error() == "no rows in result set" {
			return user, ierr.ErrNotFound
		}
		return user, err
	}

	return user, nil
}

func (u *userRepo) LookUp(ctx context.Context, id string) error {
	q := `SELECT 1 FROM users WHERE id = $1`

	v := 0
	err := u.conn.QueryRow(ctx,
		q, id).Scan(&v)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return ierr.ErrNotFound
		}
		return err
	}

	return nil
}

func (u *userRepo) UpdateAccount(ctx context.Context, id, name, url string) error {
	q := `UPDATE users SET image_url = $1, name = $2 WHERE id = $3`
	_, err := u.conn.Exec(ctx, q,
		url, name, id)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return ierr.ErrDuplicate
			}
		}
		return err
	}

	return nil
}

func (u *userRepo) GetNameBySub(ctx context.Context, id string) (string, error) {
	q := `SELECT name FROM users WHERE id = $1`

	v := ""
	err := u.conn.QueryRow(ctx,
		q, id).Scan(&v)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return "", ierr.ErrNotFound
		}
		return "", err
	}

	return v, nil
}

func (u *userRepo) GetEmailBySub(ctx context.Context, id string) (string, error) {
	q := `SELECT email FROM users WHERE id = $1`

	v := ""
	err := u.conn.QueryRow(ctx,
		q, id).Scan(&v)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return "", ierr.ErrNotFound
		}
		return "", err
	}

	return v, nil
}

// func (u *userRepo) GetNameByID(ctx context.Context, id string) (string, error) {
// 	name := ""
// 	err := u.conn.QueryRow(ctx,
// 		`SELECT name FROM users
// 		WHERE id = $1`,
// 		id).Scan(&name)
// 	if err != nil {
// 		if err.Error() == "no rows in result set" {
// 			return "", ierr.ErrNotFound
// 		}
// 		if pgErr, ok := err.(*pgconn.PgError); ok {
// 			if pgErr.Code == "22P02" {
// 				return "", ierr.ErrNotFound
// 			}
// 		}
// 		return "", err
// 	}

// 	return name, nil
// }
