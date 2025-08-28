package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"example.com/gin-memcached-base/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ItemStore interface {
	Create(ctx context.Context, name string, price float64) (models.Item, error)
	GetByID(ctx context.Context, id string) (models.Item, error)
	List(ctx context.Context) ([]models.Item, error)
}

type MongoStore struct {
	col *mongo.Collection
}

func NewMongoStore(col *mongo.Collection) *MongoStore {
	return &MongoStore{col: col}
}

func (s *MongoStore) Create(ctx context.Context, name string, price float64) (models.Item, error) {
	item := models.Item{
		ID:        uuid.NewString(),
		Name:      name,
		Price:     price,
		CreatedAt: time.Now().UTC(),
	}
	cctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err := s.col.InsertOne(cctx, item)
	return item, err
}

func (s *MongoStore) GetByID(ctx context.Context, id string) (models.Item, error) {
	cctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var out models.Item
	err := s.col.FindOne(cctx, bson.M{"_id": id}).Decode(&out)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Item{}, errors.New("not found")
		}
		return models.Item{}, err
	}
	return out, nil
}

func (s *MongoStore) List(ctx context.Context) ([]models.Item, error) {
	cctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	cur, err := s.col.Find(cctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(cctx)
	var out []models.Item
	for cur.Next(cctx) {
		var it models.Item
		if err := cur.Decode(&it); err != nil {
			return nil, err
		}
		out = append(out, it)
	}
	return out, cur.Err()
}
