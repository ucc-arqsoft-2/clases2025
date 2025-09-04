package repository

import (
    "context"
    "errors"
    "time"

    "github.com/gerasalinas/clase03-memcache-base/internal/models"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
    col *mongo.Collection
}

func NewStore(col *mongo.Collection) *Store {
    return &Store{col: col}
}

func (s *Store) List(ctx context.Context) ([]models.Item, error) {
    cur, err := s.col.Find(ctx, bson.D{})
    if err != nil {
        return nil, err
    }
    defer cur.Close(ctx)
    var out []models.Item
    for cur.Next(ctx) {
        var it models.Item
        if err := cur.Decode(&it); err != nil {
            return nil, err
        }
        out = append(out, it)
    }
    return out, cur.Err()
}

func (s *Store) Get(ctx context.Context, id primitive.ObjectID) (models.Item, error) {
    var it models.Item
    err := s.col.FindOne(ctx, bson.M{"_id": id}).Decode(&it)
    if err != nil {
        if errors.Is(err, mongo.ErrNoDocuments) {
            return it, errors.New("not found")
        }
        return it, err
    }
    return it, nil
}

func (s *Store) Create(ctx context.Context, it models.Item) (models.Item, error) {
    it.CreatedAt = time.Now()
    res, err := s.col.InsertOne(ctx, it)
    if err != nil {
        return it, err
    }
    it.ID = res.InsertedID.(primitive.ObjectID)
    return it, nil
}

func (s *Store) Update(ctx context.Context, id primitive.ObjectID, it models.Item) (models.Item, error) {
    _, err := s.col.UpdateByID(ctx, id, bson.M{"$set": bson.M{
        "name":  it.Name,
        "price": it.Price,
    }})
    if err != nil {
        return it, err
    }
    it.ID = id
    return it, nil
}

func (s *Store) Delete(ctx context.Context, id primitive.ObjectID) error {
    _, err := s.col.DeleteOne(ctx, bson.M{"_id": id})
    return err
}
