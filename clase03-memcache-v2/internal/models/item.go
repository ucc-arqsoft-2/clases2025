package models

import (
    "time"

    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Item struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name      string             `bson:"name" json:"name" binding:"required"`
    Price     float64            `bson:"price" json:"price" binding:"required"`
    CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}
