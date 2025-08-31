package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"example.com/gin-memcached-base/internal/models"
	"example.com/gin-memcached-base/internal/service"
)

type ItemHandler struct {
	service service.ItemService
}

func NewItemHandler(service service.ItemService) *ItemHandler {
	return &ItemHandler{service: service}
}

// GET /v1/items
func (h *ItemHandler) List(c *gin.Context) {
	res, err := h.service.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": res})
}

// POST /v1/items
func (h *ItemHandler) Create(c *gin.Context) {
	var in models.CreateItemInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	it, err := h.service.Create(c.Request.Context(), in.Name, in.Price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, it)
}

// GET /v1/items/:id
func (h *ItemHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	it, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"item": it})
}
