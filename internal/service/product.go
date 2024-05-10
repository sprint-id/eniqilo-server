package service

import (
	"context"
	"fmt"
	"regexp"

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
// 	"sku": "", // not null, minLength 1, maxLength 30
// 	"category": "", /** not null, enum of:
// 			- "Clothing"
// 			- "Accessories"
// 			- "Footwear"
// 			- "Beverages"
// 			*/
// 	"imageUrl": "", // not null, should be url
// 	"notes":"", // not null, minLength 1, maxLength 200
// 	"price":1, // not null, min: 1
// 	"stock": 1, // not null, min: 0, max: 100000
// 	"location": "", // not null, minLength 1, maxLength 200
// 	"isAvailable": true // not null
// }

func (u *ProductService) AddProduct(ctx context.Context, body dto.ReqAddOrUpdateProduct, sub string) (dto.ResAddOrUpdateProduct, error) {
	var res dto.ResAddOrUpdateProduct
	err := u.validator.Struct(body)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return res, ierr.ErrBadRequest
	}

	// check Image URL if invalid or not complete URL
	if !isValidURL(body.ImageUrl) {
		return res, ierr.ErrBadRequest
	}

	product := body.ToProductEntity(sub)
	res, err = u.repo.Product.AddProduct(ctx, sub, product)
	if err != nil {
		if err == ierr.ErrDuplicate {
			return res, ierr.ErrBadRequest
		}
		return res, err
	}

	return res, nil
}

func (u *ProductService) GetProduct(ctx context.Context, param dto.ParamGetProduct, sub string) ([]dto.ResGetProduct, error) {

	err := u.validator.Struct(param)
	if err != nil {
		return nil, ierr.ErrBadRequest
	}

	res, err := u.repo.Product.GetProduct(ctx, param, sub)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *ProductService) GetProductShop(ctx context.Context, param dto.ParamGetProductShop) ([]dto.ResGetProductShop, error) {
	err := u.validator.Struct(param)
	if err != nil {
		return nil, ierr.ErrBadRequest
	}

	res, err := u.repo.Product.GetProductShop(ctx, param)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *ProductService) UpdateProduct(ctx context.Context, body dto.ReqAddOrUpdateProduct, id, sub string) error {
	err := u.validator.Struct(body)
	if err != nil {
		return ierr.ErrBadRequest
	}

	// check Image URL if invalid or not complete URL
	if !isValidURL(body.ImageUrl) {
		return ierr.ErrBadRequest
	}

	product := body.ToProductEntity(sub)
	err = u.repo.Product.UpdateProduct(ctx, id, sub, product)
	if err != nil {
		return err
	}

	return nil
}

func (u *ProductService) DeleteProduct(ctx context.Context, id string, sub string) error {
	err := u.repo.Product.DeleteProduct(ctx, id, sub)
	if err != nil {
		return err
	}

	return nil
}

func isValidURL(urlString string) bool {
	// url validation using regex
	fmt.Printf("urlString: %s\n", urlString)
	regex := regexp.MustCompile(`^(https?|ftp)://[^/\s]+\.[^/\s]+(?:/.*)?(?:\.[^/\s]+)?$`)
	return regex.MatchString(urlString)
}
