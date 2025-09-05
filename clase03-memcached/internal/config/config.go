package config

import (
	"os"
	"strconv"
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
	memcachedTTL, err := strconv.Atoi(getEnv("MEMCACHED_TTL_SECONDS", "60"))
	if err != nil {
		memcachedTTL = 60
	}
	return Config{
		Port: getEnv("PORT", "8080"),
		Mongo: MongoConfig{
			URI: getEnv("MONGO_URI", "mongodb://appuser:apppass@localhost:27017/app?authSource=app"),
			DB:  getEnv("MONGO_DB", "app"),
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
