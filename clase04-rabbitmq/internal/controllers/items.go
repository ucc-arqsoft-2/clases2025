package controllers

import (
	"clase04-rabbitmq/internal/domain"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
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

// ItemsController maneja las peticiones HTTP para Items
// Responsabilidades:
// - Extraer datos del request (JSON, path params, query params)
// - Validar formato de entrada
// - Llamar al service correspondiente
// - Retornar respuesta HTTP adecuada
type ItemsController struct {
	service ItemsService // Inyección de dependencia
}

// NewItemsController crea una nueva instancia del controller
func NewItemsController(itemsService ItemsService) *ItemsController {
	return &ItemsController{
		service: itemsService,
	}
}

// GetItems maneja GET /items - Lista todos los items
// ✅ IMPLEMENTADO - Ejemplo para que los estudiantes entiendan el patrón
func (c *ItemsController) GetItems(ctx *gin.Context) {
	// 🔍 Llamar al service para obtener los datos
	items, err := c.service.List(ctx.Request.Context())
	if err != nil {
		// ❌ Error interno del servidor
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch items",
			"details": err.Error(),
		})
		return
	}

	// ✅ Respuesta exitosa con los datos
	ctx.JSON(http.StatusOK, gin.H{
		"items": items,
		"count": len(items),
	})
}

// CreateItem maneja POST /items - Crea un nuevo item
// Consigna 1: Recibir JSON, validar y crear item
func (c *ItemsController) CreateItem(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, gin.H{"error": "TODO: implementar CreateItem"})
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
