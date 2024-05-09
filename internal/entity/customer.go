package entity

type Customer struct {
	ID          string `json:"id"`
	PhoneNumber string `json:"phoneNumber"`
	Name        string `json:"name"`
	CreatedAt   string `json:"created_at"`
}
