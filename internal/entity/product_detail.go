package entity

// {
// 		"productId": "",
// 		"quantity": 1 // not null, min: 1
// 	}

type ProductDetail struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}
