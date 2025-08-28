package models

import "time"

type Item struct {
	ID        string    `json:"id" bson:"_id"`
	Name      string    `json:"name" bson:"name"`
	Price     float64   `json:"price" bson:"price"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}

type CreateItemInput struct {
	Name  string  `json:"name" binding:"required"`
	Price float64 `json:"price" binding:"required,gt=0"`
}
