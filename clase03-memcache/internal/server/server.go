package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"example.com/gin-memcached-base/internal/handlers"
	"example.com/gin-memcached-base/internal/repository"
	"example.com/gin-memcached-base/internal/service"
)

func NewRouter(store repository.ItemStore, c service.Cache) *gin.Engine {
	r := gin.Default()
	itemService := service.NewItemService(store, c)
	h := handlers.NewItemHandler(itemService)

	r.GET("/healthz", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"ok": true, "cache": "wired"})
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/items", h.List)
		v1.POST("/items", h.Create)
		v1.GET("/items/:id", h.GetByID)
	}
	return r
}
