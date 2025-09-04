package service

import (
	"context"
	"fmt"

	"example.com/gin-memcached-base/internal/models"
	"example.com/gin-memcached-base/internal/repository"
)

// Cache define la interfaz que debe implementar cualquier sistema de cache
type Cache interface {
	Get(key string) (models.Item, error)
	Set(key string, item models.Item) error

	GetList(key string) ([]models.Item, error)
	SetList(key string, items []models.Item) error
}

type ItemService interface {
	Create(ctx context.Context, name string, price float64) (models.Item, error)
	GetByID(ctx context.Context, id string) (models.Item, error)
	List(ctx context.Context) ([]models.Item, error)
}

type itemService struct {
	store repository.ItemStore
	cache Cache
}

func NewItemService(store repository.ItemStore, cache Cache) ItemService {
	return &itemService{
		store: store,
		cache: cache,
	}
}

func (s *itemService) Create(ctx context.Context, name string, price float64) (models.Item, error) {
	item, err := s.store.Create(ctx, name, price)
	if err != nil {
		return models.Item{}, err
	}

	// Guardar en cache individual
	key := fmt.Sprintf("item:%s", item.ID)
	_ = s.cache.Set(key, item)

	return item, nil
}

func (s *itemService) GetByID(ctx context.Context, id string) (models.Item, error) {
	key := fmt.Sprintf("item:%s", id)

	// Intentar primero obtener de cache
	if cached, err := s.cache.Get(key); err == nil {
		return cached, nil
	}

	// Si no está en cache, obtenerlo de la base
	item, err := s.store.GetByID(ctx, id)
	if err != nil {
		return models.Item{}, err
	}

	// Guardar en cache
	_ = s.cache.Set(key, item)

	return item, nil
}

func (s *itemService) List(ctx context.Context) ([]models.Item, error) {
	key := "items:all"

	// Intentar obtener lista del cache
	if cachedList, err := s.cache.GetList(key); err == nil {
		return cachedList, nil
	}

	// Si no está en cache, ir a la base
	items, err := s.store.List(ctx)
	if err != nil {
		return nil, err
	}

	// Guardar lista completa en cache
	_ = s.cache.SetList(key, items)

	// Opcional: cachear cada item individual también
	for _, item := range items {
		_ = s.cache.Set(fmt.Sprintf("item:%s", item.ID), item)
	}

	return items, nil
}
