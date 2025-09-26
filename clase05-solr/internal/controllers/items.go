package controllers

import (
	"clase05-solr/internal/domain"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// ItemsService define la lógica de negocio para Items
// Capa intermedia entre Controllers (HTTP) y Repository (datos)
// Responsabilidades: validaciones, transformaciones, reglas de negocio
type ItemsService interface {
	// List retorna items (si no se provee una query, retorna todos los items) de forma paginada
	List(ctx context.Context, filters domain.SearchFilters) (domain.PaginatedResponse, error)

	// Create valida y crea un nuevo item
	Create(ctx context.Context, item domain.Item) (domain.Item, error)

	// GetByID obtiene un item por su ID
	GetByID(ctx context.Context, id string) (domain.Item, error)

	// Update actualiza un item existente
	Update(ctx context.Context, id string, item domain.Item) (domain.Item, error)

	// Delete elimina un item por ID
	Delete(ctx context.Context, id string) error
}

// ItemsController maneja las peticiones HTTP para Items
// Responsabilidades:
// - Extraer datos del request (JSON, path params, query params)
// - Validar formato de entrada
// - Llamar al service correspondiente
// - Retornar respuesta HTTP adecuada
type ItemsController struct {
	service ItemsService // Inyección de dependencia
}

const (
	listDefaultPage  = 1
	listDefaultCount = 10
)

// NewItemsController crea una nueva instancia del controller
func NewItemsController(itemsService ItemsService) *ItemsController {
	return &ItemsController{
		service: itemsService,
	}
}

// List maneja GET /items - Lista los items en base a los filtros provistos en la query
func (c *ItemsController) List(ctx *gin.Context) {
	// Parsear filtros desde query params
	// Ejemplo GET /items?q=iphone&minPrice=100&maxPrice=500&page=2&count=20&sortBy=price%20desc
	filters := domain.SearchFilters{}

	filters.Name = ctx.Query("q")

	if minPriceStr := ctx.Query("minPrice"); minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			filters.MinPrice = &minPrice
		}
	}

	if maxPriceStr := ctx.Query("maxPrice"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			filters.MaxPrice = &maxPrice
		}
	}

	if pageStr := ctx.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			filters.Page = page
		}
	} else {
		filters.Page = listDefaultPage // default
	}

	if countStr := ctx.Query("count"); countStr != "" {
		if count, err := strconv.Atoi(countStr); err == nil {
			filters.Count = count
		}
	} else {
		filters.Count = listDefaultCount // default
	}

	filters.SortBy = ctx.DefaultQuery("sortBy", "createdAt desc")

	// 🔍 Llamar al service
	resp, err := c.service.List(ctx.Request.Context(), filters)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch items",
			"details": err.Error(),
		})
		return
	}

	// ✅ Respuesta exitosa con paginación incluida
	ctx.JSON(http.StatusOK, resp)
}

// CreateItem maneja POST /items - Crea un nuevo item
// Consigna 1: Recibir JSON, validar y crear item
func (c *ItemsController) CreateItem(ctx *gin.Context) {
	var item domain.Item
	if err := ctx.ShouldBindJSON(&item); err != nil {
		// ❌ Error en el formato del JSON
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid JSON format",
			"details": err.Error(),
		})
		return
	}

	created, err := c.service.Create(ctx.Request.Context(), item)
	if err != nil {
		// ❌ Error interno del servidor
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create item",
			"details": err.Error(),
		})
		return
	}

	// ✅ Respuesta exitosa con el item creado
	ctx.JSON(http.StatusCreated, gin.H{
		"item": created,
	})
}

// GetItemByID maneja GET /items/:id - Obtiene item por ID
// Consigna 2: Extraer ID del path param, validar y buscar
func (c *ItemsController) GetItemByID(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, gin.H{"error": "TODO: implementar GetItemByID"})
}

// UpdateItem maneja PUT /items/:id - Actualiza item existente
// Consigna 3: Extraer ID y datos, validar y actualizar
func (c *ItemsController) UpdateItem(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, gin.H{"error": "TODO: implementar UpdateItem"})
}

// DeleteItem maneja DELETE /items/:id - Elimina item por ID
// Consigna 4: Extraer ID, validar y eliminar
func (c *ItemsController) DeleteItem(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, gin.H{"error": "TODO: implementar DeleteItem"})
}

// 📚 Notas sobre HTTP Status Codes
//
// 200 OK - Operación exitosa con contenido
// 201 Created - Recurso creado exitosamente
// 204 No Content - Operación exitosa sin contenido (típico para DELETE)
// 400 Bad Request - Error en los datos enviados por el cliente
// 404 Not Found - Recurso no encontrado
// 500 Internal Server Error - Error interno del servidor
// 501 Not Implemented - Funcionalidad no implementada (para TODOs)
//
// 💡 Tip: En una API real, sería buena práctica crear una función
// helper para manejar respuestas de error de manera consistente
