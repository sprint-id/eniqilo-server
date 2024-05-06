package handler

import (
	"encoding/json"
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

func (h *productHandler) AddCat(w http.ResponseWriter, r *http.Request) {
	var req dto.ReqAddOrUpdateCat
	var res dto.ResAddCat
	var jsonData map[string]interface{}

	// Decode request body into the jsonData map
	err := json.NewDecoder(r.Body).Decode(&jsonData)
	if err != nil {
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}

	// Check for unexpected fields
	expectedFields := []string{"name", "race", "sex", "ageInMonth", "description", "imageUrls"}
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

	// Validate "race" value
	if req.Race == "Maine Coon" {
		req.Race = "Maine_Coon"
	} else if req.Race == "British Shorthair" {
		req.Race = "British_Shorthair"
	} else if req.Race == "Scottish Fold" {
		req.Race = "Scottish_Fold"
	}

	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}

	res, err = h.productSvc.AddCat(r.Context(), req, token.Subject())
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

func (h *productHandler) GetCat(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	var param dto.ParamGetCat

	param.ID = queryParams.Get("id")
	param.Limit, _ = strconv.Atoi(queryParams.Get("limit"))
	param.Offset, _ = strconv.Atoi(queryParams.Get("offset"))
	param.Race = queryParams.Get("race")
	param.Sex = queryParams.Get("sex")
	param.HasMatched, _ = strconv.ParseBool(queryParams.Get("hasMatched"))
	param.AgeInMonth = queryParams.Get("ageInMonth")
	param.Owned, _ = strconv.ParseBool(queryParams.Get("owned"))
	param.Search = queryParams.Get("search")

	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}

	cats, err := h.productSvc.GetCat(r.Context(), param, token.Subject())
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	successRes := response.SuccessReponse{}
	successRes.Message = "success"
	successRes.Data = cats

	json.NewEncoder(w).Encode(successRes)
	w.WriteHeader(http.StatusOK)
}

func (h *productHandler) GetCatByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}

	cat, err := h.productSvc.GetCatByID(r.Context(), id, token.Subject())
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	successRes := response.SuccessReponse{}
	successRes.Message = "success"
	successRes.Data = cat

	json.NewEncoder(w).Encode(successRes)
	w.WriteHeader(http.StatusOK)
}

func (h *productHandler) UpdateCat(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req dto.ReqAddOrUpdateCat
	var jsonData map[string]interface{}

	// Decode request body into the jsonData map
	err := json.NewDecoder(r.Body).Decode(&jsonData)
	if err != nil {
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}

	// Check for unexpected fields
	expectedFields := []string{"name", "race", "sex", "ageInMonth", "description", "imageUrls"}
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

	// Validate "race" value
	if req.Race == "Maine Coon" {
		req.Race = "Maine_Coon"
	} else if req.Race == "British Shorthair" {
		req.Race = "British_Shorthair"
	} else if req.Race == "Scottish Fold" {
		req.Race = "Scottish_Fold"
	}

	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}

	err = h.productSvc.UpdateCat(r.Context(), req, id, token.Subject())
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *productHandler) DeleteCat(w http.ResponseWriter, r *http.Request) {
	// Get id from URL path parameters
	id := r.PathValue("id")
	// fmt.Printf("id: %s\n", id)

	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}

	if id == "" {
		http.Error(w, "id is required for cat", http.StatusBadRequest)
		return
	}

	err = h.productSvc.DeleteCat(r.Context(), id, token.Subject())
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
