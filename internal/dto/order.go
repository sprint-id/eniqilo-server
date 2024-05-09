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
	ReqOrder struct {
		CustomerID     string             `json:"customerId" validate:"required"`
		ProductDetails []ReqProductDetail `json:"productDetails" validate:"required,min=1,dive"`
		Paid           int                `json:"paid" validate:"required,min=1"`
		Change         int                `json:"change" validate:"required,min=0"`
	}

	ReqProductDetail struct {
		ProductID string `json:"productId" validate:"required"`
		Quantity  int    `json:"quantity" validate:"required,min=1"`
	}

	ResOrder struct {
		ID        string `json:"id"`
		CreatedAt string `json:"createdAt"`
	}
)

// to entity order
func (r *ReqOrder) ToOrderEntity() *entity.Order {
	var productDetails []entity.ProductDetail
	for _, pd := range r.ProductDetails {
		productDetails = append(productDetails, entity.ProductDetail{
			ProductID: pd.ProductID,
			Quantity:  pd.Quantity,
		})
	}

	return &entity.Order{
		CustomerID:     r.CustomerID,
		ProductDetails: productDetails,
		Paid:           r.Paid,
		Change:         r.Change,
	}
}
