package dao

import (
	"clase04-rabbitmq/internal/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Item struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	Price     float64            `bson:"price"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

// ToDomain convierte de modelo DB a modelo de negocio
func (d Item) ToDomain() domain.Item {
	return domain.Item{
		ID:        d.ID.Hex(), // ObjectID -> string
		Name:      d.Name,
		Price:     d.Price,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

// FromDomain convierte de modelo de negocio a modelo DB
func FromDomain(domainItem domain.Item) Item {
	// Si el ID está vacío, DB generará uno automáticamente
	var objectID primitive.ObjectID
	if domainItem.ID != "" {
		objectID, _ = primitive.ObjectIDFromHex(domainItem.ID)
	}

	return Item{
		ID:        objectID,
		Name:      domainItem.Name,
		Price:     domainItem.Price,
		CreatedAt: domainItem.CreatedAt,
		UpdatedAt: domainItem.UpdatedAt,
	}
}
