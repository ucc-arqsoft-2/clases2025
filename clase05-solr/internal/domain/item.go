package domain

import (
	"time"
)

type SearchFilters struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	MinPrice *float64 `json:"min_price"`
	MaxPrice *float64 `json:"max_price"`
	SortBy   string   `json:"sort_by"`
	Page     int      `json:"page"`
	Count    int      `json:"count"`
}

type PaginatedResponse struct {
	Page    int    `json:"page"`
	Count   int    `json:"count"`
	Total   int    `json:"total"`
	Results []Item `json:"results"`
}

type Item struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
