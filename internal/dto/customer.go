package dto

import "github.com/sprint-id/eniqilo-server/internal/entity"

// {
// 	"phoneNumber": "+628123123123", // not null, minLength: 10, maxLength: 16, should start with `+` and international calling codes
// 	// reference: https://countrycode.org
// 	// it should support country code like `591` and `1-246` as well
// 	// customer phoneNumber shoud be a different entity from staff phoneNumber
// 	"name": "namadepan namabelakang" // not null, minLength 5, maxLength 50, name can be duplicate with others
// }

type (
	ReqRegisterCustomer struct {
		PhoneNumber string `json:"phoneNumber" validate:"required,min=10,max=16"`
		Name        string `json:"name" validate:"required,min=5,max=50"`
	}

	ResRegisterOrGetCustomer struct {
		UserID      string `json:"userId"`
		PhoneNumber string `json:"phoneNumber"`
		Name        string `json:"name"`
	}

	ParamGetCustomer struct {
		PhoneNumber string `json:"phoneNumber"`
		Name        string `json:"name"`
	}
)

// ToEntity to convert dto to entity
func (d *ReqRegisterCustomer) ToCustomerEntity() entity.Customer {
	return entity.Customer{
		PhoneNumber: d.PhoneNumber,
		Name:        d.Name,
	}
}
