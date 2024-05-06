package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type CredentialType string

const (
	PhoneType CredentialType = "phone"
	EmailType CredentialType = "email"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func New() *validator.Validate {
	cv := &CustomValidator{
		Validator: validator.New(),
	}

	return cv.Validator
}

func (cv *CustomValidator) Validate(i interface{}) error {
	err := cv.Validator.Struct(i)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldError := range validationErrors {
			err := fmt.Errorf("%s is %s", fieldError.Field(), fieldError.Tag())
			return err
		}
	}

	return nil
}
