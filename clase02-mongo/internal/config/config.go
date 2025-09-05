package config

import (
	"os"
)

type Config struct {
	Port  string
	Mongo MongoConfig
}

type MongoConfig struct {
	URI string
	DB  string
}

func Load() Config {
	return Config{
		Port: getEnv("PORT", "8080"),
		Mongo: MongoConfig{
			URI: getEnv("MONGO_URI", "mongodb://appuser:apppass@localhost:27017/app?authSource=app"),
			DB:  getEnv("MONGO_DB", "app"),
		},
	}
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
