package service

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/sprint-id/eniqilo-server/internal/cfg"
	"github.com/sprint-id/eniqilo-server/internal/dto"
	"github.com/sprint-id/eniqilo-server/internal/ierr"
	"github.com/sprint-id/eniqilo-server/internal/repo"
)

type TransactionService struct {
	repo      *repo.Repo
	validator *validator.Validate
	cfg       *cfg.Cfg
}

func newTransactionService(repo *repo.Repo, validator *validator.Validate, cfg *cfg.Cfg) *TransactionService {
	return &TransactionService{repo, validator, cfg}
}

func (ts *TransactionService) AddTransaction(ctx context.Context, body dto.ReqTransaction, sub string) (dto.ResTransaction, error) {
	var res dto.ResTransaction
	err := ts.validator.Struct(body)
	if err != nil {
		return res, ierr.ErrBadRequest
	}

	transaction := body.ToTransactionEntity()
	res, err = ts.repo.Transaction.AddTransaction(ctx, sub, *transaction)
	if err != nil {
		if err == ierr.ErrDuplicate {
			return res, ierr.ErrDuplicate
		}
		return res, err
	}

	return res, nil
}

func (ts *TransactionService) GetTransactionHistory(ctx context.Context, param dto.ParamGetTransactionHistory, sub string) ([]dto.ResOrderHistory, error) {
	var res []dto.ResOrderHistory
	err := ts.validator.Struct(param)
	if err != nil {
		return res, ierr.ErrBadRequest
	}

	res, err = ts.repo.Transaction.GetTransactionHistory(ctx, param, sub)
	if err != nil {
		return res, err
	}

	return res, nil
}
