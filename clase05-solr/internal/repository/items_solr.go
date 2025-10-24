package repository

import (
	"clase05-solr/internal/clients"
	"clase05-solr/internal/domain"
	"context"
	"fmt"
	"strings"
)

type SolrClient interface {
}

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

// List retorna items desde Solr en base a los filtros
func (r *SolrItemsRepository) List(ctx context.Context, filters domain.SearchFilters) (domain.PaginatedResponse, error) {
	query := buildQuery(filters)
	return r.client.Search(ctx, query, filters.Page, filters.Count)
}

// GetByID busca un item por su ID en Solr
func (r *SolrItemsRepository) GetByID(ctx context.Context, id string) (domain.Item, error) {
	results, err := r.List(ctx, domain.SearchFilters{ID: id})
	if err != nil {
		return domain.Item{}, fmt.Errorf("error searching item by ID in solr: %w", err)
	}
	if results.Total == 0 {
		return domain.Item{}, fmt.Errorf("item with ID %s not found", id)
	}
	return results.Results[0], nil
}

// Create indexa un nuevo item en Solr
func (r *SolrItemsRepository) Create(ctx context.Context, item domain.Item) (domain.Item, error) {
	if err := r.client.Index(ctx, item); err != nil {
		return domain.Item{}, fmt.Errorf("error indexing item in solr: %w", err)
	}
	return item, nil
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

// buildQuery construye la query de Solr a partir de los filtros
func buildQuery(filters domain.SearchFilters) string {
	var parts []string

	// Si no hay filtros, devolvemos todo
	if filters.ID == "" && filters.Name == "" && filters.MinPrice == nil && filters.MaxPrice == nil {
		return "*:*"
	}

	// Filtro por ID
	if filters.ID != "" {
		parts = append(parts, fmt.Sprintf("id:%s", filters.ID))
	}

	// Filtro por nombre
	if filters.Name != "" {
		parts = append(parts, fmt.Sprintf("name:*%s*", filters.Name))
	}

	// Filtro por rango de precios
	if filters.MinPrice != nil && filters.MaxPrice != nil {
		parts = append(parts, fmt.Sprintf("price:[%f TO %f]", *filters.MinPrice, *filters.MaxPrice))
	}

	return strings.Join(parts, " AND ")
}
