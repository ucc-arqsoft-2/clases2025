package main

import (
	"context"
	"log"
	"os"
	"time"

	"example.com/gin-memcached-base/internal/repository"
	"example.com/gin-memcached-base/internal/server"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	addr := getEnv("ADDR", ":8080")
	// TODO(Clase): Leer MEMCACHED_ADDR y CACHE_TTL_SECONDS de variables de entorno
	// memcachedAddr := getEnv("MEMCACHED_ADDR", "memcached:11211")
	// ttl := getEnv("CACHE_TTL_SECONDS", "60")

	mongoURI := getEnv("MONGO_URI", "mongodb://mongo:27017")
	mongoDB := getEnv("MONGO_DB", "demo")
	mongoColl := getEnv("MONGO_COLLECTION", "items")

	// 1) Conexión Mongo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("mongo connect: %v", err)
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("mongo ping: %v", err)
	}
	col := client.Database(mongoDB).Collection(mongoColl)
	store := repository.NewMongoStore(col)

	// 2) TODO(Clase): Crear cliente de Memcached con TTL (usar internal/cache/memcached.go)
	// cacheClient := cache.New(memcachedAddr, ttlDur)

	// 3) Router (por ahora, sin cache - solo store)
	r := server.NewRouter(store)

	log.Printf("listening on %s | mongo=%s/%s", addr, mongoDB, mongoColl)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}