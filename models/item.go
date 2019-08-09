package models

// Item Structure
type Item struct {
	Name            string  `json:"name"`
	AlternativeName string  `json:"alternative_name"`
	Price           float64 `json:"price"`
	CategoryID      int     `json:"category_id"`
}
