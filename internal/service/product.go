package service

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/sprint-id/eniqilo-server/internal/cfg"
	"github.com/sprint-id/eniqilo-server/internal/dto"
	"github.com/sprint-id/eniqilo-server/internal/ierr"
	"github.com/sprint-id/eniqilo-server/internal/repo"
)

type ProductService struct {
	repo      *repo.Repo
	validator *validator.Validate
	cfg       *cfg.Cfg
}

func newProductService(repo *repo.Repo, validator *validator.Validate, cfg *cfg.Cfg) *ProductService {
	return &ProductService{repo, validator, cfg}
}

// {
// 	"name": "", // not null, minLength 1, maxLength 30
// 	"race": "", /** not null, enum of:
// 			- "Persian"
// 			- "Maine Coon"
// 			- "Siamese"
// 			- "Ragdoll"
// 			- "Bengal"
// 			- "Sphynx"
// 			- "British Shorthair"
// 			- "Abyssinian"
// 			- "Scottish Fold"
// 			- "Birman" */
// 	"sex": "", // not null, enum of: "male" / "female"
// 	"ageInMonth": 1, // not null, min: 1, max: 120082
// 	"description":"" // not null, minLength 1, maxLength 200
// 	"imageUrls":[ // not null, minItems: 1, items: not null, should be url
// 		"","",""
// 	]
// }

func (u *ProductService) AddCat(ctx context.Context, body dto.ReqAddOrUpdateCat, sub string) (dto.ResAddCat, error) {
	var res dto.ResAddCat
	err := u.validator.Struct(body)
	if err != nil {
		return res, ierr.ErrBadRequest
	}

	cat := body.ToCatEntity(sub)
	res, err = u.repo.Product.AddCat(ctx, sub, cat)
	if err != nil {
		if err == ierr.ErrDuplicate {
			return res, ierr.ErrBadRequest
		}
		return res, err
	}

	return res, nil
}

func (u *ProductService) GetCat(ctx context.Context, param dto.ParamGetCat, sub string) ([]dto.ResGetCat, error) {

	err := u.validator.Struct(param)
	if err != nil {
		return nil, ierr.ErrBadRequest
	}

	res, err := u.repo.Product.GetCat(ctx, param, sub)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *ProductService) GetCatByID(ctx context.Context, id, sub string) (dto.ResGetCat, error) {
	res, err := u.repo.Product.GetCatByID(ctx, id, sub)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (u *ProductService) UpdateCat(ctx context.Context, body dto.ReqAddOrUpdateCat, id, sub string) error {
	err := u.validator.Struct(body)
	if err != nil {
		return ierr.ErrBadRequest
	}

	cat := body.ToCatEntity(sub)
	err = u.repo.Product.UpdateCat(ctx, id, sub, cat)
	if err != nil {
		return err
	}

	return nil
}

func (u *ProductService) DeleteCat(ctx context.Context, id string, sub string) error {
	err := u.repo.Product.DeleteCat(ctx, id, sub)
	if err != nil {
		return err
	}

	return nil
}
