package service

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/gerasalinas/clase03-memcache-base/internal/cache"
    "github.com/gerasalinas/clase03-memcache-base/internal/models"
    "github.com/gerasalinas/clase03-memcache-base/internal/repository"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// Ejercicios en esta capa:
// - Cachear el listado (clave 'items:all') con TTL 60s
// - Cachear el item individual (clave 'item:<id>')
// - Invalidar 'items:all' y 'item:<id>' al crear, actualizar o eliminar

type ItemsService struct {
    repo *repository.Store
    c    cache.Cache
}

func NewItemsService(repo *repository.Store, c cache.Cache) *ItemsService {
    return &ItemsService{repo: repo, c: c}
}

func (s *ItemsService) List(ctx context.Context) ([]models.Item, error) {
    // TODO: intentar leer desde cache
    out, err := s.repo.List(ctx)
    if err != nil {
        return nil, err
    }
    // TODO: escribir en cache con TTL 60s
    _ = json.Marshal // para que el import no quede sin usar
    return out, nil
}

func (s *ItemsService) Get(ctx context.Context, hex string) (models.Item, error) {
    id, err := primitive.ObjectIDFromHex(hex)
    if err != nil { return models.Item{}, fmt.Errorf("invalid id: %w", err) }

    // TODO: intentar leer desde cache 'item:<id>'
    it, err := s.repo.Get(ctx, id)
    if err != nil { return it, err }
    // TODO: escribir en cache con TTL 60s
    _ = time.Second // para que el import no quede sin usar
    return it, nil
}

func (s *ItemsService) Create(ctx context.Context, in models.Item) (models.Item, error) {
    it, err := s.repo.Create(ctx, in)
    if err != nil { return it, err }
    // TODO: invalidar 'items:all' y 'item:<id>'
    return it, nil
}

func (s *ItemsService) Update(ctx context.Context, hex string, in models.Item) (models.Item, error) {
    id, err := primitive.ObjectIDFromHex(hex)
    if err != nil { return models.Item{}, fmt.Errorf("invalid id: %w", err) }
    it, err := s.repo.Update(ctx, id, in)
    if err != nil { return it, err }
    // TODO: invalidar 'items:all' y 'item:<id>'
    return it, nil
}

func (s *ItemsService) Delete(ctx context.Context, hex string) error {
    id, err := primitive.ObjectIDFromHex(hex)
    if err != nil { return fmt.Errorf("invalid id: %w", err) }
    if err := s.repo.Delete(ctx, id); err != nil { return err }
    // TODO: invalidar 'items:all' y 'item:<id>'
    return nil
}
