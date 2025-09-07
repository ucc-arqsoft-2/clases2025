package services

import (
	"clase02-mongo/internal/domain"
	"context"
	"errors"
	"fmt"
	"strings"
)

// ItemsRepository define las operaciones de datos para Items
// Patrón Repository: abstrae el acceso a datos del resto de la aplicación
type ItemsRepository interface {
	// List retorna todos los items de la base de datos
	List(ctx context.Context) ([]domain.Item, error)

	// Create inserta un nuevo item en DB
	Create(ctx context.Context, item domain.Item) (domain.Item, error)

	// GetByID busca un item por su ID
	GetByID(ctx context.Context, id string) (domain.Item, error)

	// Update actualiza un item existente
	Update(ctx context.Context, id string, item domain.Item) (domain.Item, error)

	// Delete elimina un item por ID
	Delete(ctx context.Context, id string) error
} // ItemsServiceImpl implementa ItemsService

type ItemsServiceImpl struct {
	repository ItemsRepository // Inyección de dependencia
}

// NewItemsService crea una nueva instancia del service
// Pattern: Dependency Injection - recibe dependencies como parámetros
func NewItemsService(repository ItemsRepository) ItemsServiceImpl {
	return ItemsServiceImpl{
		repository: repository,
	}
}

// List obtiene todos los items
// ✅ IMPLEMENTADO - Delegación simple al repository
func (s *ItemsServiceImpl) List(ctx context.Context) ([]domain.Item, error) {
	// En este caso, no hay lógica de negocio especial
	// Solo delegamos al repository
	return s.repository.List(ctx)
}

// Create valida y crea un nuevo item
// Consigna 1: Validar name no vacío y price >= 0
func (s *ItemsServiceImpl) Create(ctx context.Context, item domain.Item) (domain.Item, error) {
	// Validar campos del item
	if err := s.validateItem(item); err != nil {
		return domain.Item{}, fmt.Errorf("validation error: %w", err)
	}

	created, err := s.repository.Create(ctx, item)
	if err != nil {
		return domain.Item{}, fmt.Errorf("error creating item in repository: %w", err)
	}

	return created, nil
}

// GetByID obtiene un item por su ID
// Consigna 2: Validar formato de ID antes de consultar DB
func (s *ItemsServiceImpl) GetByID(ctx context.Context, id string) (domain.Item, error) {
	item, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return domain.Item{}, fmt.Errorf("error getting item from repository: %w", err)
	}

	return item, nil
}

// Update actualiza un item existente
// Consigna 3: Validar campos antes de actualizar
func (s *ItemsServiceImpl) Update(ctx context.Context, id string, item domain.Item) (domain.Item, error) {
	_, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return domain.Item{}, fmt.Errorf("item does not exists: %w", err)
	}

	// Validar campos del item
	if err := s.validateItem(item); err != nil {
		return domain.Item{}, fmt.Errorf("validation error: %w", err)
	}

	// Actualizar en la BD
	updated, err := s.repository.Update(ctx, id, item)
	if err != nil {
		return domain.Item{}, fmt.Errorf("error updating item in repository: %w", err)
	}

	return updated, nil
}

// Delete elimina un item por ID
// Consigna 4: Validar ID antes de eliminar
func (s *ItemsServiceImpl) Delete(ctx context.Context, id string) error {
	_, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("item does not exists: %w", err)
	}

	err = s.repository.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("error deleting item in repository: %w", err)
	}

	return nil
}

// validateItem aplica reglas de negocio para validar un item
// 🎯 Función helper para reutilizar validaciones
func (s *ItemsServiceImpl) validateItem(item domain.Item) error {
	// 📝 Name es obligatorio y no puede estar vacío
	if strings.TrimSpace(item.Name) == "" {
		return errors.New("name is required and cannot be empty")
	}

	// 💰 Price debe ser >= 0 (productos gratis están permitidos)
	if item.Price < 0 {
		return errors.New("price must be greater than or equal to 0")
	}

	// ✅ Todas las validaciones pasaron
	return nil
}
