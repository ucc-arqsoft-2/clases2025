package main

import (
	"clase02-mongo/internal/config"
	"clase02-mongo/internal/controllers"
	"clase02-mongo/internal/db"
	"clase02-mongo/internal/middleware"
	"clase02-mongo/internal/repository"
	"clase02-mongo/internal/services"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// 📋 Cargar configuración desde las variables de entorno
	cfg := config.Load()

	// 🔌 Conectar a MongoDB
	// Establecemos conexión y obtenemos referencia a la base de datos
	_, mongoDB, err := db.Connect(context.Background(), cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		log.Fatalf("mongo connect error: %v", err)
	}

	// 🏗️ Inicializar capas de la aplicación (Dependency Injection)
	// Patrón: Repository -> Service -> Controller
	// Cada capa tiene una responsabilidad específica

	// Capa de datos: maneja operaciones MongoDB
	itemRepo := repository.NewMongoItemsRepository(mongoDB)

	// Capa de lógica de negocio: validaciones, transformaciones
	itemService := services.NewItemsService(&itemRepo)

	// Capa de controladores: maneja HTTP requests/responses
	itemController := controllers.NewItemsController(&itemService)

	// 🌐 Configurar router HTTP con Gin
	router := gin.Default()

	// Middleware: funciones que se ejecutan en cada request
	router.Use(middleware.CORSMiddleware)

	// 🏥 Health check endpoint
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 📚 Rutas de Items API
	// GET /items - listar todos los items (✅ implementado)
	router.GET("/items", itemController.GetItems)

	// TODO: Implementar la lógica de estos endpoints (actualmente retornan 501)
	// POST /items - crear nuevo item
	router.POST("/items", itemController.CreateItem)

	// GET /items/:id - obtener item por ID
	router.GET("/items/:id", itemController.GetItemByID)

	// PUT /items/:id - actualizar item existente
	router.PUT("/items/:id", itemController.UpdateItem)

	// DELETE /items/:id - eliminar item
	router.DELETE("/items/:id", itemController.DeleteItem)

	// Configuración del server HTTP
	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	log.Printf("🚀 API listening on port %s", cfg.Port)
	log.Printf("📊 Health check: http://localhost:%s/healthz", cfg.Port)
	log.Printf("📚 Items API: http://localhost:%s/items", cfg.Port)

	// Iniciar servidor (bloquea hasta que se pare el servidor)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
