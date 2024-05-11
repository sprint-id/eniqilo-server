package dto

import "github.com/sprint-id/eniqilo-server/internal/entity"

// {
// 	"customerId": "", // ID Should be string
// 	"productDetails": [
// 		{
// 			"productId": "",
// 			"quantity": 1 // not null, min: 1
// 		}
// 	], // ID Should be string, minItems: 1
// 	"paid": 1, // not null, min: 1, validate the change based on all product price
// 	"change": 0, // not null, min 0
// }

type (
	ReqTransaction struct {
		CustomerID     string             `json:"customerId" validate:"required"`
		ProductDetails []ReqProductDetail `json:"productDetails" validate:"required,min=1,dive"`
		Paid           int                `json:"paid" validate:"required,min=1"`
		Change         *int               `json:"change" validate:"required,min=0"`
	}

	ReqProductDetail struct {
		ProductID string `json:"productId" validate:"required"`
		Quantity  int    `json:"quantity" validate:"required,min=1"`
	}

	ParamGetTransactionHistory struct {
		CustomerID string `json:"customerId"`
		Limit      int    `json:"limit"`
		Offset     int    `json:"offset"`
		CreatedAt  string `json:"createdAt"`
	}

	ResTransaction struct {
		ID        string `json:"id"`
		CreatedAt string `json:"createdAt"`
	}

	ResTransactionHistory struct {
		TransactionID  string             `json:"transactionId"`
		CustomerID     string             `json:"customerId"`
		ProductDetails []ReqProductDetail `json:"productDetails"`
		Paid           int                `json:"paid"`
		Change         int                `json:"change"`
		CreatedAt      string             `json:"createdAt"`
	}
)

// to entity transaction
func (r *ReqTransaction) ToTransactionEntity() *entity.Transaction {
	var productDetails []entity.ProductDetail
	for _, pd := range r.ProductDetails {
		productDetails = append(productDetails, entity.ProductDetail{
			ProductID: pd.ProductID,
			Quantity:  pd.Quantity,
		})
	}

	return &entity.Transaction{
		CustomerID:     r.CustomerID,
		ProductDetails: productDetails,
		Paid:           r.Paid,
		Change:         *r.Change,
	}
}
