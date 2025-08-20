package services

import (
	"clase02-mongo/internal/domain"
	"clase02-mongo/internal/repository"
	"context"
	"errors"
	"strings"
)

// ItemsService define la lógica de negocio para Items
// Capa intermedia entre Controllers (HTTP) y Repository (datos)
// Responsabilidades: validaciones, transformaciones, reglas de negocio
type ItemsService interface {
	// List retorna todos los items (sin filtros por ahora)
	List(ctx context.Context) ([]domain.Item, error)
	
	// Create valida y crea un nuevo item
	Create(ctx context.Context, item domain.Item) (domain.Item, error)
	
	// GetByID obtiene un item por su ID
	GetByID(ctx context.Context, id string) (domain.Item, error)
	
	// Update actualiza un item existente
	Update(ctx context.Context, id string, item domain.Item) (domain.Item, error)
	
	// Delete elimina un item por ID
	Delete(ctx context.Context, id string) error
}

// ItemsServiceImpl implementa ItemsService
type ItemsServiceImpl struct {
	repository repository.ItemsRepository // Inyección de dependencia
}

// NewItemsService crea una nueva instancia del service
// Pattern: Dependency Injection - recibe dependencies como parámetros
func NewItemsService(repository repository.ItemsRepository) ItemsService {
	return &ItemsServiceImpl{repository: repository}
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
	return domain.Item{}, errors.New("TODO: implementar Create")
}

// GetByID obtiene un item por su ID
// Consigna 2: Validar formato de ID antes de consultar DB
func (s *ItemsServiceImpl) GetByID(ctx context.Context, id string) (domain.Item, error) {
	return domain.Item{}, errors.New("TODO: implementar GetByID")
}

// Update actualiza un item existente
// Consigna 3: Validar campos antes de actualizar
func (s *ItemsServiceImpl) Update(ctx context.Context, id string, item domain.Item) (domain.Item, error) {
	return domain.Item{}, errors.New("TODO: implementar Update")
}

// Delete elimina un item por ID
// Consigna 4: Validar ID antes de eliminar
func (s *ItemsServiceImpl) Delete(ctx context.Context, id string) error {
	return errors.New("TODO: implementar Delete")
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
