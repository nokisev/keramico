package models

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Rating      int     `json:"rating"`
	Price       float64 `json:"price"`
}
