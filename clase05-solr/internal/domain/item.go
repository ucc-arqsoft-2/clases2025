package domain

import (
	"time"
)

type SearchResponse struct {
	Page    int    `json:"page"`
	Count   int    `json:"count"`
	Results []Item `json:"results"`
}

type Item struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
