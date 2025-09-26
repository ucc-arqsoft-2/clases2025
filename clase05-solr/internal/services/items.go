package services

import (
	"clase05-solr/internal/domain"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
)

// ItemsRepository define las operaciones de datos para Items
// Patrón Repository: abstrae el acceso a datos del resto de la aplicación
type ItemsRepository interface {
	// List retorna items de la base de datos en base a los filtros
	List(ctx context.Context, filters domain.SearchFilters) (domain.PaginatedResponse, error)

	// Create inserta un nuevo item en DB
	Create(ctx context.Context, item domain.Item) (domain.Item, error)

	// GetByID busca un item por su ID
	GetByID(ctx context.Context, id string) (domain.Item, error)

	// Update actualiza un item existente
	Update(ctx context.Context, id string, item domain.Item) (domain.Item, error)

	// Delete elimina un item por ID
	Delete(ctx context.Context, id string) error
} // ItemsServiceImpl implementa ItemsService

type ItemsPublisher interface {
	Publish(ctx context.Context, action string, itemID string) error
}

type ItemsConsumer interface {
	Consume(ctx context.Context, handler func(ctx context.Context, message ItemEvent) error) error
}

type ItemsServiceImpl struct {
	repository ItemsRepository // Inyección de dependencia
	cache      ItemsRepository // Inyección de dependencia
	search     ItemsRepository // Repositorio de búsqueda (Solr)
	publisher  ItemsPublisher
	consumer   ItemsConsumer
}

// NewItemsService crea una nueva instancia del service
// Pattern: Dependency Injection - recibe dependencies como parámetros
func NewItemsService(repository ItemsRepository, cache ItemsRepository, search ItemsRepository, publisher ItemsPublisher, consumer ItemsConsumer) ItemsServiceImpl {
	return ItemsServiceImpl{
		repository: repository,
		cache:      cache,
		search:     search,
		publisher:  publisher,
		consumer:   consumer,
	}
}

// List obtiene todos los items
// ✅ IMPLEMENTADO - Delegación simple al repository
func (s *ItemsServiceImpl) List(ctx context.Context, filters domain.SearchFilters) (domain.PaginatedResponse, error) {
	// En este caso, no hay lógica de negocio especial
	// Solo delegamos al search repository
	return s.search.List(ctx, filters)
}

// Create valida y crea un nuevo item
// Consigna 1: Validar name no vacío y price >= 0
func (s *ItemsServiceImpl) Create(ctx context.Context, item domain.Item) (domain.Item, error) {
	created, err := s.repository.Create(ctx, item)
	if err != nil {
		return domain.Item{}, fmt.Errorf("error creating item in repository: %w", err)
	}

	if err := s.publisher.Publish(ctx, "create", created.ID); err != nil {
		return domain.Item{}, fmt.Errorf("error publishing item creation: %w", err)
	}

	_, err = s.cache.Create(ctx, created)
	if err != nil {
		return domain.Item{}, fmt.Errorf("error creating item in cache: %w", err)
	}

	// TODO: Publicar a RabbitMQ

	return created, nil
}

// GetByID obtiene un item por su ID
// Consigna 2: Validar formato de ID antes de consultar DB
func (s *ItemsServiceImpl) GetByID(ctx context.Context, id string) (domain.Item, error) {
	item, err := s.cache.GetByID(ctx, id)
	if err != nil {
		item, err := s.repository.GetByID(ctx, id)
		if err != nil {
			return domain.Item{}, fmt.Errorf("error getting item from repository: %w", err)
		}

		_, err = s.cache.Create(ctx, item)
		if err != nil {
			return domain.Item{}, fmt.Errorf("error creating item in cache: %w", err)
		}

		return item, nil
	}
	return item, nil
}

// Update actualiza un item existente
// Consigna 3: Validar campos antes de actualizar
func (s *ItemsServiceImpl) Update(ctx context.Context, id string, item domain.Item) (domain.Item, error) {

	// TODO: Actualizar en DB
	// TODO: Guardar en cache
	// TODO: Publicar a RabbitMQ

	return domain.Item{}, errors.New("TODO: implementar Update")
}

// Delete elimina un item por ID
// Consigna 4: Validar ID antes de eliminar
func (s *ItemsServiceImpl) Delete(ctx context.Context, id string) error {

	// TODO: Borrar de cache
	// TODO: Borrar de DB
	// TODO: Publicar a RabbitMQ

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

type ItemEvent struct {
	Action string `json:"action"` // "create", "update", "delete"
	ItemID string `json:"item_id"`
}

func (s *ItemsServiceImpl) InitConsumer(ctx context.Context) {
	// Iniciar Go routine para el consumer
	slog.Info("🐰 Starting RabbitMQ consumer...")

	if err := s.consumer.Consume(ctx, s.handleMessage); err != nil {
		slog.Error("❌ Error in RabbitMQ consumer: %v", err)
	}
	slog.Info("🐰 RabbitMQ consumer stopped.")
}

// handleMessage procesa los mensajes recibidos de RabbitMQ
func (s *ItemsServiceImpl) handleMessage(ctx context.Context, message ItemEvent) error {
	slog.Info("📨 Processing message",
		slog.String("action", message.Action),
		slog.String("item_id", message.ItemID),
	)

	switch message.Action {
	case "create":
		slog.Info("✅ Item created", slog.String("item_id", message.ItemID))

		// Indexar el item en Solr para búsquedas
		item, err := s.repository.GetByID(ctx, message.ItemID)
		if err != nil {
			slog.Error("❌ Error getting item for indexing",
				slog.String("item_id", message.ItemID),
				slog.String("error", err.Error()))
			return fmt.Errorf("error getting item for indexing: %w", err)
		}

		if _, err := s.search.Create(ctx, item); err != nil {
			slog.Error("❌ Error indexing item in search",
				slog.String("item_id", message.ItemID),
				slog.String("error", err.Error()))
			return fmt.Errorf("error indexing item in search: %w", err)
		}

		slog.Info("🔍 Item indexed in search engine", slog.String("item_id", message.ItemID))

	case "update":
		slog.Info("✏️ Item updated", slog.String("item_id", message.ItemID))
		// Invalidar cache si es necesario

	case "delete":
		slog.Info("🗑️ Item deleted", slog.String("item_id", message.ItemID))
		// Limpiar cache, logs, etc.

	default:
		slog.Info("⚠️ Unknown action", slog.String("action", message.Action))
	}

	return nil
}
