package service

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/sprint-id/eniqilo-server/internal/cfg"
	"github.com/sprint-id/eniqilo-server/internal/dto"
	"github.com/sprint-id/eniqilo-server/internal/ierr"
	"github.com/sprint-id/eniqilo-server/internal/repo"
)

type OrderService struct {
	repo      *repo.Repo
	validator *validator.Validate
	cfg       *cfg.Cfg
}

func newOrderService(repo *repo.Repo, validator *validator.Validate, cfg *cfg.Cfg) *OrderService {
	return &OrderService{repo, validator, cfg}
}

func (o *OrderService) AddOrder(ctx context.Context, body dto.ReqOrder, sub string) (dto.ResOrder, error) {
	var res dto.ResOrder
	err := o.validator.Struct(body)
	if err != nil {
		return res, ierr.ErrBadRequest
	}

	order := body.ToOrderEntity()
	res, err = o.repo.Order.AddOrder(ctx, sub, *order)
	if err != nil {
		if err == ierr.ErrDuplicate {
			return res, ierr.ErrBadRequest
		}
		return res, err
	}

	return res, nil
}
