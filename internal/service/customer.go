package service

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/sprint-id/eniqilo-server/internal/cfg"
	"github.com/sprint-id/eniqilo-server/internal/dto"
	"github.com/sprint-id/eniqilo-server/internal/ierr"
	"github.com/sprint-id/eniqilo-server/internal/repo"
)

type CustomerService struct {
	repo      *repo.Repo
	validator *validator.Validate
	cfg       *cfg.Cfg
}

func newCustomerService(repo *repo.Repo, validator *validator.Validate, cfg *cfg.Cfg) *CustomerService {
	return &CustomerService{repo, validator, cfg}
}

func (u *CustomerService) RegisterCustomer(ctx context.Context, body dto.ReqRegisterCustomer, sub string) (dto.ResRegisterOrGetCustomer, error) {
	var res dto.ResRegisterOrGetCustomer
	err := u.validator.Struct(body)
	if err != nil {
		return res, ierr.ErrBadRequest
	}

	customer := body.ToCustomerEntity()
	res, err = u.repo.Customer.RegisterCustomer(ctx, sub, customer)
	if err != nil {
		if err == ierr.ErrDuplicate {
			return res, ierr.ErrBadRequest
		}
		return res, err
	}

	return res, nil
}

func (u *CustomerService) GetCustomer(ctx context.Context, param dto.ParamGetCustomer, sub string) ([]dto.ResRegisterOrGetCustomer, error) {
	res, err := u.repo.Customer.GetCustomer(ctx, param, sub)
	if err != nil {
		return nil, err
	}

	return res, nil
}
