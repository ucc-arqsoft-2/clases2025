package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port      string
	Mongo     MongoConfig
	Memcached MemcachedConfig
	RabbitMQ  RabbitMQConfig
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

type RabbitMQConfig struct {
	Username  string
	Password  string
	QueueName string
	Host      string
	Port      string
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
		RabbitMQ: RabbitMQConfig{
			Username:  getEnv("RABBITMQ_USERNAME", "guest"),
			Password:  getEnv("RABBITMQ_PASSWORD", "guest"),
			QueueName: getEnv("RABBITMQ_QUEUE_NAME", "items-news"),
			Host:      getEnv("RABBITMQ_HOST", "localhost"),
			Port:      getEnv("RABBITMQ_PORT", "5672"),
		},
	}
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
