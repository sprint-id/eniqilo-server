package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/sprint-id/eniqilo-server/internal/dto"
	"github.com/sprint-id/eniqilo-server/internal/ierr"
	"github.com/sprint-id/eniqilo-server/internal/service"
	response "github.com/sprint-id/eniqilo-server/pkg/resp"
)

type customerHandler struct {
	customerSvc *service.CustomerService
}

func newCustomerHandler(customerSvc *service.CustomerService) *customerHandler {
	return &customerHandler{customerSvc}
}

// {
// 	"phoneNumber": "+628123123123", // not null, minLength: 10, maxLength: 16, should start with `+` and international calling codes
// 	// reference: https://countrycode.org
// 	// it should support country code like `591` and `1-246` as well
// 	// customer phoneNumber shoud be a different entity from staff phoneNumber
// 	"name": "namadepan namabelakang" // not null, minLength 5, maxLength 50, name can be duplicate with others
// }

func (h *customerHandler) RegisterCustomer(w http.ResponseWriter, r *http.Request) {
	var req dto.ReqRegisterCustomer
	var res dto.ResRegisterOrGetCustomer
	var jsonData map[string]interface{}

	// Decode request body into the jsonData map
	err := json.NewDecoder(r.Body).Decode(&jsonData)
	if err != nil {
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}

	// Check for unexpected fields
	expectedFields := []string{"phoneNumber", "name"}
	for key := range jsonData {
		if !contains(expectedFields, key) {
			http.Error(w, "unexpected field in request body: "+key, http.StatusBadRequest)
			return
		}
	}

	// Convert the jsonData map into the req struct
	bytes, err := json.Marshal(jsonData)
	if err != nil {
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(bytes, &req)
	if err != nil {
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}

	// show request
	fmt.Printf("RegisterCustomer request: %+v\n", req)

	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}

	res, err = h.customerSvc.RegisterCustomer(r.Context(), req, token.Subject())
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	successRes := response.SuccessReponse{}
	successRes.Message = "success"
	successRes.Data = res

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // Set HTTP status code to 201
	err = json.NewEncoder(w).Encode(successRes)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *customerHandler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	var param dto.ParamGetCustomer

	param.PhoneNumber = queryParams.Get("phoneNumber")
	param.Name = queryParams.Get("name")

	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}

	customers, err := h.customerSvc.GetCustomer(r.Context(), param, token.Subject())
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	// show response
	// fmt.Printf("GetMatch response: %+v\n", customers)

	successRes := response.SuccessReponse{}
	successRes.Message = "success"
	successRes.Data = customers

	json.NewEncoder(w).Encode(successRes)
	w.WriteHeader(http.StatusOK)
}
