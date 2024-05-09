package dto

import (
	"github.com/sprint-id/eniqilo-server/internal/entity"
)

type (
	ReqAddOrUpdateProduct struct {
		Name        string `json:"name" validate:"required,min=1,max=30"`
		SKU         string `json:"sku" validate:"required,min=1,max=30"`
		Category    string `json:"category" validate:"required,oneof=Clothing Accessories Footwear Beverages"`
		ImageUrl    string `json:"imageUrl" validate:"required,url"`
		Notes       string `json:"notes" validate:"required,min=1,max=200"`
		Price       int    `json:"price" validate:"required,min=1"`
		Stock       int    `json:"stock" validate:"required,min=0,max=100000"`
		Location    string `json:"location" validate:"required,min=1,max=200"`
		IsAvailable bool   `json:"isAvailable" validate:"required"`
	}

	ParamGetProduct struct {
		ID          string `json:"id"`
		Limit       int    `json:"limit"`
		Offset      int    `json:"offset"`
		Name        string `json:"name"`
		IsAvailable string `json:"isAvailable"`
		Category    string `json:"category"`
		SKU         string `json:"sku"`
		Price       string `json:"price"`
		InStock     string `json:"inStock"`
		CreatedAt   string `json:"createdAt"`
		Search      string `json:"search"`
	}

	ParamGetProductShop struct {
		Limit    int    `json:"limit"`
		Offset   int    `json:"offset"`
		Name     string `json:"name"`
		Category string `json:"category"`
		SKU      string `json:"sku"`
		Price    string `json:"price"`
		InStock  bool   `json:"inStock"`
	}

	ResAddOrUpdateProduct struct {
		ID        string `json:"id"`
		CreatedAt string `json:"createdAt"`
	}

	ResGetProduct struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		SKU         string `json:"sku"`
		Category    string `json:"category"`
		ImageUrl    string `json:"imageUrl"`
		Stock       int    `json:"stock"`
		Notes       string `json:"notes"`
		Price       int    `json:"price"`
		Location    string `json:"location"`
		IsAvailable bool   `json:"isAvailable"`
		CreatedAt   string `json:"createdAt"`
	}

	ResGetProductShop struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		SKU       string `json:"sku"`
		Category  string `json:"category"`
		ImageUrl  string `json:"imageUrl"`
		Stock     int    `json:"stock"`
		Price     int    `json:"price"`
		Location  string `json:"location"`
		CreatedAt string `json:"createdAt"`
	}
)

func (d *ReqAddOrUpdateProduct) ToProductEntity(userId string) entity.Product {
	return entity.Product{
		Name:        d.Name,
		SKU:         d.SKU,
		Category:    d.Category,
		ImageUrl:    d.ImageUrl,
		Notes:       d.Notes,
		Price:       d.Price,
		Stock:       d.Stock,
		Location:    d.Location,
		IsAvailable: d.IsAvailable,
		UserID:      userId,
	}
}
