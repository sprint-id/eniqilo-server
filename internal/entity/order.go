package entity

// Order struct
type Order struct {
	CustomerID     string          `json:"customerId"`
	ProductDetails []ProductDetail `json:"productDetails"`
	Paid           int             `json:"paid"`
	Change         int             `json:"change"`
}
