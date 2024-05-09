package entity

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

type Product struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	SKU         string `json:"sku"`
	Category    string `json:"category"`
	ImageUrl    string `json:"imageUrl"`
	Notes       string `json:"notes"`
	Price       int    `json:"price"`
	Stock       int    `json:"stock"`
	Location    string `json:"location"`
	IsAvailable bool   `json:"isAvailable"`
	CreatedAt   string `json:"created_at"`

	UserID string `json:"user_id"`
}
