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

type matchRepo struct {
	conn *pgxpool.Pool
}

func newMatchRepo(conn *pgxpool.Pool) *matchRepo {
	return &matchRepo{conn}
}

func (mr *matchRepo) MatchCat(ctx context.Context, sub string, match_cat entity.MatchCat) (dto.ResMatchCat, error) {
	q := `INSERT INTO match_cats (user_id, match_cat_id, user_cat_id, message, created_at)
	VALUES ( $1, $2, $3, $4, EXTRACT(EPOCH FROM now())::bigint) RETURNING id`

	// show the query
	// fmt.Println(q)

	var id string
	err := mr.conn.QueryRow(ctx, q, sub, match_cat.MatchCatId, match_cat.UserCatId, match_cat.Message).Scan(&id)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return dto.ResMatchCat{}, ierr.ErrDuplicate
			}
		}
		return dto.ResMatchCat{}, err
	}

	return dto.ResMatchCat{MatchId: id}, nil

	// _, err := mr.conn.Exec(ctx, q,
	// 	sub, match_cat.MatchCatId, match_cat.UserCatId, match_cat.Message)

	// if err != nil {
	// 	if pgErr, ok := err.(*pgconn.PgError); ok {
	// 		if pgErr.Code == "23505" {
	// 			return ierr.ErrDuplicate
	// 		}
	// 	}
	// 	return err
	// }

	// return nil
}

func (mr *matchRepo) GetMatch(ctx context.Context, sub string) ([]dto.ResGetMatchCat, error) {
	q := `SELECT mc.id, u.name, u.email, u.created_at, 
		mcd.id, mcd.name, mcd.race, mcd.sex, mcd.description, mcd.age_in_month, mcd.image_urls, mcd.has_matched, mcd.created_at, 
		ucd.id, ucd.name, ucd.race, ucd.sex, ucd.description, ucd.age_in_month, ucd.image_urls, ucd.has_matched, ucd.created_at,
		mc.message, mc.created_at
		FROM match_cats mc
		INNER JOIN users u ON mc.user_id = u.id
		INNER JOIN cats mcd ON mc.user_cat_id = mcd.id
		INNER JOIN cats ucd ON mc.user_cat_id = ucd.id
		WHERE 1=1
		ORDER BY mc.created_at DESC`

	rows, err := mr.conn.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []dto.ResGetMatchCat
	for rows.Next() {
		var match dto.ResGetMatchCat
		var issuedByCreatedAt int64
		var matchCatCreatedAt int64
		var userCatCreatedAt int64
		var matchCreatedAt int64
		err = rows.Scan(&match.ID, &match.IssuedBy.Name, &match.IssuedBy.Email, &issuedByCreatedAt,
			&match.MatchCatDetail.ID, &match.MatchCatDetail.Name, &match.MatchCatDetail.Race, &match.MatchCatDetail.Sex, &match.MatchCatDetail.Description, &match.MatchCatDetail.AgeInMonth, &match.MatchCatDetail.ImageUrls, &match.MatchCatDetail.HasMatched, &matchCatCreatedAt,
			&match.UserCatDetail.ID, &match.UserCatDetail.Name, &match.UserCatDetail.Race, &match.UserCatDetail.Sex, &match.UserCatDetail.Description, &match.UserCatDetail.AgeInMonth, &match.UserCatDetail.ImageUrls, &match.UserCatDetail.HasMatched, &userCatCreatedAt,
			&match.Message, &matchCreatedAt)
		if err != nil {
			return nil, err
		}

		match.IssuedBy.CreatedAt = timepkg.TimeToISO8601(time.Unix(issuedByCreatedAt, 0))
		match.MatchCatDetail.CreatedAt = timepkg.TimeToISO8601(time.Unix(matchCatCreatedAt, 0))
		match.UserCatDetail.CreatedAt = timepkg.TimeToISO8601(time.Unix(userCatCreatedAt, 0))
		match.CreatedAt = timepkg.TimeToISO8601(time.Unix(matchCreatedAt, 0))
		matches = append(matches, match)
	}

	return matches, nil
}

func (mr *matchRepo) ApproveMatch(ctx context.Context, sub string, match_id string) error {
	approveq := `UPDATE match_cats SET has_approved = true WHERE user_id = $1 AND id = $2`
	mcq := `UPDATE cats c SET has_matched = true FROM match_cats mc WHERE mc.user_id = $1 AND mc.id = $2 AND c.id = mc.match_cat_id`
	ucq := `UPDATE cats c SET has_matched = true FROM match_cats mc WHERE mc.user_id = $1 AND mc.id = $2 AND c.id = mc.user_cat_id`

	_, err := mr.conn.Exec(ctx, approveq, sub, match_id)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return ierr.ErrNotFound
			}
		}
		return err
	}

	_, err = mr.conn.Exec(ctx, mcq, sub, match_id)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return ierr.ErrNotFound
			}
		}
		return err
	}

	_, err = mr.conn.Exec(ctx, ucq, sub, match_id)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return ierr.ErrNotFound
			}
		}
		return err
	}

	// get match_cat_id and user_cat_id from match_id
	var match_cat_id, user_cat_id string
	gmcucq := `SELECT match_cat_id, user_cat_id FROM match_cats WHERE id = $1`
	err = mr.conn.QueryRow(ctx, gmcucq, match_id).Scan(&match_cat_id, &user_cat_id)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return ierr.ErrNotFound
			}
		}
		return err
	}

	// delete all match_id that contains match_cat_id or user_cat_id except the current match_id
	_, err = mr.conn.Exec(ctx, `DELETE FROM match_cats WHERE (match_cat_id = $1 OR match_cat_id = $2 OR user_cat_id = $1 OR user_cat_id = $2) AND id != $3`, match_cat_id, user_cat_id, match_id)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return ierr.ErrNotFound
			}
		}
		return err
	}

	return nil
}

func (mr *matchRepo) RejectMatch(ctx context.Context, sub string, match_id string) error {
	q := `UPDATE match_cats SET has_approved = false WHERE user_id = $1 AND id = $2`

	_, err := mr.conn.Exec(ctx, q, sub, match_id)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return ierr.ErrNotFound
			}
		}
		return err
	}

	return nil
}

func (mr *matchRepo) DeleteMatch(ctx context.Context, sub string, match_id string) error {
	q := `DELETE FROM match_cats WHERE user_id = $1 AND id = $2`

	_, err := mr.conn.Exec(ctx, q, sub, match_id)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return ierr.ErrNotFound
			}
		}
		return err
	}

	return nil
}
