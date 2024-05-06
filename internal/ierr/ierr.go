package ierr

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

type customError struct {
	Message string `json:"message"`
}

func (e customError) Error() string {
	return e.Message
}

func ExtendErr(err customError, msg string) error {
	err.Message = fmt.Sprintf("%s, err : %s", err.Message, msg)
	return err
}

var (
	ErrInternal   = customError{Message: "Sorry, an internal server error occurred. Please try again later."}
	ErrDuplicate  = customError{Message: "The data you provided conflicts with existing data. Please review the information you entered"}
	ErrNotFound   = customError{Message: "Sorry, the resource you requested could not be found."}
	ErrBadRequest = customError{Message: "Sorry, the request is invalid. Please check your input and try again."}
	ErrForbidden  = customError{Message: "You do not have permission to access or edit this resource."}
)

func TranslateError(err error) (code int, msg string) {
	log.Println(err)
	
	switch errors.Cause(err) {
	case ErrDuplicate:
		return http.StatusConflict, err.Error()
	case ErrNotFound:
		return http.StatusNotFound, err.Error()
	case ErrForbidden:
		return http.StatusForbidden, err.Error()
	case ErrBadRequest:
		return http.StatusBadRequest, err.Error()
	}

	return http.StatusInternalServerError, ErrInternal.Message
}
