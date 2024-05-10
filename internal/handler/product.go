package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/jwtauth/v5"
	"github.com/sprint-id/eniqilo-server/internal/dto"
	"github.com/sprint-id/eniqilo-server/internal/ierr"
	"github.com/sprint-id/eniqilo-server/internal/service"
	response "github.com/sprint-id/eniqilo-server/pkg/resp"
)

type productHandler struct {
	productSvc *service.ProductService
}

func newProductHandler(productSvc *service.ProductService) *productHandler {
	return &productHandler{productSvc}
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

func (h *productHandler) AddProduct(w http.ResponseWriter, r *http.Request) {
	var req dto.ReqAddOrUpdateProduct
	var res dto.ResAddOrUpdateProduct
	var jsonData map[string]interface{}

	// Decode request body into the jsonData map
	err := json.NewDecoder(r.Body).Decode(&jsonData)
	if err != nil {
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}

	// Check if the payload is empty
	if len(jsonData) == 0 {
		http.Error(w, "empty payload", http.StatusUnauthorized)
		return
	}

	// Check for unexpected fields
	expectedFields := []string{"name", "sku", "category", "imageUrl", "notes", "price", "stock", "location", "isAvailable"}
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

	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}

	res, err = h.productSvc.AddProduct(r.Context(), req, token.Subject())
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

func (h *productHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	var param dto.ParamGetProduct

	param.ID = queryParams.Get("id")
	param.Limit, _ = strconv.Atoi(queryParams.Get("limit"))
	param.Offset, _ = strconv.Atoi(queryParams.Get("offset"))
	param.Name = queryParams.Get("name")
	param.IsAvailable = queryParams.Get("isAvailable")
	param.Category = queryParams.Get("category")
	param.SKU = queryParams.Get("sku")
	param.Price = queryParams.Get("price")
	param.InStock = queryParams.Get("inStock")
	param.CreatedAt = queryParams.Get("createdAt")

	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}

	products, err := h.productSvc.GetProduct(r.Context(), param, token.Subject())
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	successRes := response.SuccessReponse{}
	successRes.Message = "success"
	successRes.Data = products

	json.NewEncoder(w).Encode(successRes)
	w.WriteHeader(http.StatusOK)
}

func (h *productHandler) GetProductShop(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	var param dto.ParamGetProductShop

	param.Limit, _ = strconv.Atoi(queryParams.Get("limit"))
	param.Offset, _ = strconv.Atoi(queryParams.Get("offset"))
	param.Name = queryParams.Get("name")
	param.Category = queryParams.Get("category")
	param.SKU = queryParams.Get("sku")
	param.Price = queryParams.Get("price")
	param.InStock = queryParams.Get("inStock")

	products, err := h.productSvc.GetProductShop(r.Context(), param)
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	successRes := response.SuccessReponse{}
	successRes.Message = "success"
	successRes.Data = products

	json.NewEncoder(w).Encode(successRes)
	w.WriteHeader(http.StatusOK)
}

func (h *productHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req dto.ReqAddOrUpdateProduct
	var jsonData map[string]interface{}

	// Decode request body into the jsonData map
	err := json.NewDecoder(r.Body).Decode(&jsonData)
	if err != nil {
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}

	// id should be a number
	_, err = strconv.Atoi(id)
	if err != nil {
		http.Error(w, "id should be a number", http.StatusNotFound)
		return
	}

	// Check if the payload is empty
	if len(jsonData) == 0 || id == "" {
		http.Error(w, "empty payload", http.StatusUnauthorized)
		return
	}

	// Check for unexpected fields
	expectedFields := []string{"name", "sku", "category", "imageUrl", "notes", "price", "stock", "location", "isAvailable"}
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

	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}

	err = h.productSvc.UpdateProduct(r.Context(), req, id, token.Subject())
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *productHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	// Get id from URL path parameters
	id := r.PathValue("id")
	fmt.Printf("id: %s\n", id)

	// id should be a number
	_, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "id should be a number", http.StatusNotFound)
		return
	}

	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}

	if id == "" {
		http.Error(w, "id is required for cat", http.StatusBadRequest)
		return
	}

	err = h.productSvc.DeleteProduct(r.Context(), id, token.Subject())
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// The contains function checks if a slice contains a string
func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}
