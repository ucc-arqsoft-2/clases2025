package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	// "example.com/gin-memcached-base/internal/cache" // TODO(Clase)
	"example.com/gin-memcached-base/internal/handlers"
	"example.com/gin-memcached-base/internal/repository"
)

func NewRouter(store repository.ItemStore /*, c *cache.Client*/) *gin.Engine {
	r := gin.Default()
	h := handlers.NewItemHandler(store /*, c*/)

	// TODO(Clase): Cuando implementen cache, usar c.SelfTest aquí.
	r.GET("/healthz", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"ok": true, "cache": "not-wired-yet"})
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/items", h.List)
		v1.POST("/items", h.Create)
		v1.GET("/items/:id", h.GetByID)
	}
	return r
}
