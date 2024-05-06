package dto

import (
	"github.com/sprint-id/eniqilo-server/internal/entity"
	"github.com/sprint-id/eniqilo-server/pkg/auth"
)

type (
	ReqRegister struct {
		PhoneNumber string `json:"phoneNumber" validate:"required,min=10,max=16"`
		Name        string `json:"name" validate:"required,min=5,max=50"`
		Password    string `json:"password" validate:"required,min=5,max=15"`
	}
	ResRegister struct {
		UserID      string `json:"userId"`
		PhoneNumber string `json:"phoneNumber,omitempty"`
		Name        string `json:"name"`
		AccessToken string `json:"accessToken"`
	}
	ReqLogin struct {
		PhoneNumber string `json:"phoneNumber" validate:"required,min=10,max=16"`
		Password    string `json:"password" validate:"required,min=5,max=15"`
	}
	ResLogin struct {
		UserID      string `json:"userId"`
		PhoneNumber string `json:"phoneNumber,omitempty"`
		Name        string `json:"name"`
		AccessToken string `json:"accessToken"`
	}
	ReqUpdateAccount struct {
		ImageURL string `json:"imageUrl" validate:"required,url"`
		Name     string `json:"name" validate:"required,min=5,max=50"`
	}
)

func (d *ReqRegister) ToEntity(cryptCost int) entity.User {
	return entity.User{Name: d.Name, Password: auth.HashPassword(d.Password, cryptCost), PhoneNumber: d.PhoneNumber}
}
