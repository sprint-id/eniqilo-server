package entity

type Cat struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Race        string   `json:"race"`
	Sex         string   `json:"sex"`
	AgeInMonth  int      `json:"age_in_month"`
	Description string   `json:"description"`
	ImageUrls   []string `json:"image_urls"`
	HasMatched  bool     `json:"has_matched"`
	CreatedAt   string   `json:"created_at"`

	UserID string `json:"user_id"`
}
