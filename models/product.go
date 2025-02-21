package models

import "time"

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	SellerID    int       `json:"seller_id"`
	CreatedAt   time.Time `json:"created_at"`
}
