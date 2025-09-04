package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/gerasalinas/clase03-memcache-base/internal/cache"
    "github.com/gerasalinas/clase03-memcache-base/internal/models"
    "github.com/gerasalinas/clase03-memcache-base/internal/service"
)

func RegisterRoutes(r *gin.Engine, svc *service.ItemsService, c cache.Cache) {
    r.GET("/healthz", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"status": "ok"}) })

    r.GET("/items", func(ctx *gin.Context) {
        items, err := svc.List(ctx)
        if err != nil { ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }
        ctx.JSON(http.StatusOK, gin.H{"items": items})
    })

    r.GET("/items/:id", func(ctx *gin.Context) {
        it, err := svc.Get(ctx, ctx.Param("id"))
        if err != nil { ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()}); return }
        ctx.JSON(http.StatusOK, it)
    })

    r.POST("/items", func(ctx *gin.Context) {
        var in models.Item
        if err := ctx.ShouldBindJSON(&in); err != nil {
            ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
        }
        it, err := svc.Create(ctx, in)
        if err != nil { ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }
        ctx.JSON(http.StatusCreated, it)
    })

    r.PUT("/items/:id", func(ctx *gin.Context) {
        var in models.Item
        if err := ctx.ShouldBindJSON(&in); err != nil {
            ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
        }
        it, err := svc.Update(ctx, ctx.Param("id"), in)
        if err != nil { ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }
        ctx.JSON(http.StatusOK, it)
    })

    r.DELETE("/items/:id", func(ctx *gin.Context) {
        if err := svc.Delete(ctx, ctx.Param("id")); err != nil {
            ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
        }
        ctx.Status(http.StatusNoContent)
    })

    // Endpoints para inspeccionar la cache (ejercicio: hacer que funcionen)
    r.GET("/__cache/keys", func(ctx *gin.Context) {
        keys, _ := c.Keys()
        ctx.JSON(http.StatusOK, gin.H{"keys": keys})
    })
    r.GET("/__cache/get", func(ctx *gin.Context) {
        key := ctx.Query("key")
        if key == "" { ctx.JSON(http.StatusBadRequest, gin.H{"error":"missing key"}); return }
        bs, ok := c.Get(key)
        if !ok { ctx.JSON(http.StatusNotFound, gin.H{"error":"cache miss"}); return }
        ctx.Data(http.StatusOK, "application/json; charset=utf-8", bs)
    })
}
