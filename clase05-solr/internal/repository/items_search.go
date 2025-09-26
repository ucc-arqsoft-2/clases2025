package repository

import (
	"clase05-solr/internal/clients"
	"clase05-solr/internal/domain"
	"context"
	"fmt"
)

// SolrItemsRepository implementa ItemsRepository usando Solr
type SolrItemsRepository struct {
	client *clients.SolrClient
}

// NewSolrItemsRepository crea una nueva instancia del repository
func NewSolrItemsRepository(host, port, core string) *SolrItemsRepository {
	client := clients.NewSolrClient(host, port, core)
	return &SolrItemsRepository{
		client: client,
	}
}

// List retorna todos los items indexados en Solr
func (r *SolrItemsRepository) List(ctx context.Context) (domain.SearchResponse, error) {
	items, err := r.client.GetAll(ctx)
	if err != nil {
		return domain.SearchResponse{}, fmt.Errorf("error listing items from solr: %w", err)
	}
	return domain.SearchResponse{
		Page:    1,
		Count:   len(items),
		Results: items,
	}, nil
}

// Create indexa un nuevo item en Solr
func (r *SolrItemsRepository) Create(ctx context.Context, item domain.Item) (domain.Item, error) {
	if err := r.client.Index(ctx, item); err != nil {
		return domain.Item{}, fmt.Errorf("error indexing item in solr: %w", err)
	}
	return item, nil
}

// GetByID busca un item por su ID en Solr
func (r *SolrItemsRepository) GetByID(ctx context.Context, id string) (domain.Item, error) {
	query := fmt.Sprintf("id:%s", id)
	items, err := r.client.Search(ctx, query)
	if err != nil {
		return domain.Item{}, fmt.Errorf("error searching item in solr: %w", err)
	}

	if len(items) == 0 {
		return domain.Item{}, fmt.Errorf("item with id %s not found", id)
	}

	return items[0], nil
}

// Update actualiza un item existente en Solr (re-indexa)
func (r *SolrItemsRepository) Update(ctx context.Context, id string, item domain.Item) (domain.Item, error) {
	// En Solr, actualizar es equivalente a re-indexar con el mismo ID
	item.ID = id
	if err := r.client.Index(ctx, item); err != nil {
		return domain.Item{}, fmt.Errorf("error updating item in solr: %w", err)
	}
	return item, nil
}

// Delete elimina un item por ID de Solr
func (r *SolrItemsRepository) Delete(ctx context.Context, id string) error {
	if err := r.client.Delete(ctx, id); err != nil {
		return fmt.Errorf("error deleting item from solr: %w", err)
	}
	return nil
}

// Search busca items por query personalizada
func (r *SolrItemsRepository) Search(ctx context.Context, query string) ([]domain.Item, error) {
	return r.client.Search(ctx, query)
}

// SearchByName busca items por nombre
func (r *SolrItemsRepository) SearchByName(ctx context.Context, name string) ([]domain.Item, error) {
	query := fmt.Sprintf("name:*%s*", name)
	return r.client.Search(ctx, query)
}

// SearchByPriceRange busca items en un rango de precios
func (r *SolrItemsRepository) SearchByPriceRange(ctx context.Context, minPrice, maxPrice float64) ([]domain.Item, error) {
	query := fmt.Sprintf("price:[%f TO %f]", minPrice, maxPrice)
	return r.client.Search(ctx, query)
}
