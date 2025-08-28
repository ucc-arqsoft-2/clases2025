package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"example.com/gin-memcached-base/internal/models"
	"example.com/gin-memcached-base/internal/repository"
	// "example.com/gin-memcached-base/internal/cache" // TODO(Clase)
)

const cacheListKey = "items:all"

type ItemHandler struct {
	store repository.ItemStore
	// cache *cache.Client // TODO(Clase)
}

func NewItemHandler(store repository.ItemStore /*, c *cache.Client*/) *ItemHandler {
	// return &ItemHandler{store: store, cache: c} // TODO(Clase)
	return &ItemHandler{store: store}
}

// GET /v1/items
func (h *ItemHandler) List(c *gin.Context) {
	// TODO(Clase): Implementar cache-first usando h.cache
	res, err := h.store.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"source": "store", "items": res})
}

// POST /v1/items
func (h *ItemHandler) Create(c *gin.Context) {
	var in models.CreateItemInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	it, err := h.store.Create(c.Request.Context(), in.Name, in.Price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// TODO(Clase):
	// - Setear cache para "item:<id>"
	// - Invalidar cache de lista "items:all"
	c.JSON(http.StatusCreated, it)
}

// GET /v1/items/:id
func (h *ItemHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	// TODO(Clase): Intentar cache-first "item:<id>"
	it, err := h.store.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	// TODO(Clase): Si viene de store, setear cache
	c.JSON(http.StatusOK, gin.H{"source": "store", "item": it})
}
