package service

import (
	"context"

	"example.com/gin-memcached-base/internal/models"
	"example.com/gin-memcached-base/internal/repository"
)

// Cache define la interfaz que debe implementar cualquier sistema de cache
// Principio: "accept interfaces, return types" - el service define lo que necesita
type Cache interface {
	Get(key string) (models.Item, error)
	Set(key string, item models.Item) error
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

	// TODO(Clase): Setear el item en cache usando s.cache.Set()

	return item, nil
}

func (s *itemService) GetByID(ctx context.Context, id string) (models.Item, error) {
	// TODO(Clase): Intentar obtener el item del cache primero usando s.cache.Get()
	// Si existe en cache, retornarlo directamente

	// Si no está en cache, obtenerlo de la base de datos
	item, err := s.store.GetByID(ctx, id)
	if err != nil {
		return models.Item{}, err
	}

	// TODO(Clase): Guardar el item en cache para futuras consultas usando s.cache.Set()
	
	return item, nil
}

func (s *itemService) List(ctx context.Context) ([]models.Item, error) {
	// TODO(Clase - Opcional): Implementar cache para la lista completa
	// Pueden usar una key como "items:all" para cachear la lista
	return s.store.List(ctx)
}