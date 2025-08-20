package config

import "os"

type Config struct {
	Port     string
	MongoURI string
	MongoDB  string
}

func Load() Config {
	return Config{
		Port:     getEnv("PORT", "8080"),
		MongoURI: getEnv("MONGODB_URI", "mongodb://appuser:apppass@localhost:27017/app?authSource=app"),
		MongoDB:  getEnv("MONGODB_DB", "app"),
	}
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
