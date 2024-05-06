package entity

type User struct {
	ID          string `json:"id"`
	PhoneNumber string `json:"phone_number"`
	Name        string `json:"name"`
	Password    string `json:"password"`
	CreatedAt   string `json:"created_at"` // TODO accshualllyy, we dont need this
}
