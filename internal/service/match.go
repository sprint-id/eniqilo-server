package service

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/sprint-id/eniqilo-server/internal/cfg"
	"github.com/sprint-id/eniqilo-server/internal/dto"
	"github.com/sprint-id/eniqilo-server/internal/ierr"
	"github.com/sprint-id/eniqilo-server/internal/repo"
)

type MatchService struct {
	repo      *repo.Repo
	validator *validator.Validate
	cfg       *cfg.Cfg
}

func newMatchService(repo *repo.Repo, validator *validator.Validate, cfg *cfg.Cfg) *MatchService {
	return &MatchService{repo, validator, cfg}
}

func (u *MatchService) MatchCat(ctx context.Context, body dto.ReqMatchCat, sub string) (dto.ResMatchCat, error) {
	var res dto.ResMatchCat
	err := u.validator.Struct(body)
	if err != nil {
		return res, ierr.ErrBadRequest
	}

	_, err = u.repo.Product.GetCatByID(ctx, body.UserCatId, sub)
	if err != nil {
		fmt.Println("error get cat by id")
		return res, ierr.ErrBadRequest
	}

	// get name from sub
	name, err := u.repo.User.GetNameBySub(ctx, sub)
	// fmt.Println(name)
	if err != nil {
		fmt.Println("error get name by sub")
		return res, err
	}

	// get email from sub
	email, err := u.repo.User.GetEmailBySub(ctx, sub)
	// fmt.Println(email)
	if err != nil {
		fmt.Println("error get email by sub")
		return res, err
	}

	match := body.ToMatchCatEntity(name, email)
	res, err = u.repo.Match.MatchCat(ctx, sub, match)
	if err != nil {
		if err == ierr.ErrDuplicate {
			return res, ierr.ErrBadRequest
		}
		return res, err
	}

	return res, nil
}

func (u *MatchService) GetMatch(ctx context.Context, sub string) ([]dto.ResGetMatchCat, error) {
	res, err := u.repo.Match.GetMatch(ctx, sub)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *MatchService) ApproveMatch(ctx context.Context, body dto.ReqApproveOrRejectMatchCat, sub string) error {
	err := u.validator.Struct(body)
	if err != nil {
		return ierr.ErrBadRequest
	}

	err = u.repo.Match.ApproveMatch(ctx, sub, body.MatchId)
	if err != nil {
		if err == ierr.ErrNotFound {
			return ierr.ErrNotFound
		}
		return err
	}

	return nil
}

func (u *MatchService) RejectMatch(ctx context.Context, body dto.ReqApproveOrRejectMatchCat, sub string) error {
	err := u.validator.Struct(body)
	if err != nil {
		return ierr.ErrBadRequest
	}

	err = u.repo.Match.RejectMatch(ctx, sub, body.MatchId)
	if err != nil {
		if err == ierr.ErrNotFound {
			return ierr.ErrNotFound
		}
		return err
	}

	return nil
}

func (u *MatchService) DeleteMatch(ctx context.Context, sub string, matchId string) error {
	err := u.repo.Match.DeleteMatch(ctx, sub, matchId)
	if err != nil {
		if err == ierr.ErrNotFound {
			return ierr.ErrNotFound
		}
		return err
	}

	return nil
}
