package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	Mongo     MongoConfig
	Memcached MemcachedConfig
}

type MongoConfig struct {
	URI string
	DB  string
}

type MemcachedConfig struct {
	Host       string
	Port       string
	TTLSeconds int
}

func Load() Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env file")
	}

	memcachedTTL, err := strconv.Atoi(getEnv("MEMCACHED_TTL_SECONDS", "60"))
	if err != nil {
		memcachedTTL = 60
	}
	return Config{
		Port: getEnv("PORT", "8080"),
		Mongo: MongoConfig{
			URI: getEnv("MONGO_URI", "mongodb://localhost:27017"),
			DB:  getEnv("MONGO_DB", "demo"),
		},
		Memcached: MemcachedConfig{
			Host:       getEnv("MEMCACHED_HOST", "localhost"),
			Port:       getEnv("MEMCACHED_PORT", "11211"),
			TTLSeconds: memcachedTTL,
		},
	}
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
