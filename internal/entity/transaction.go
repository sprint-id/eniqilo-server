package entity

// Transaction struct
type Transaction struct {
	CustomerID     string          `json:"customerId"`
	ProductDetails []ProductDetail `json:"productDetails"`
	Paid           int             `json:"paid"`
	Change         int             `json:"change"`
}
