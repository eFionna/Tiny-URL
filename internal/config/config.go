package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	BaseURL       string
	RedisAddr     string
	RedisPassword string
	RedisDB       int

	RedisPoolSize     int
	RedisMinIdleConns int
}

func Load() *Config {
	if os.Getenv("GO_ENV") != "production" {
		_ = godotenv.Load()
	}

	return &Config{
		Port:          mustGet("PORT"),
		BaseURL:       mustGet("BASE_URL"),
		RedisAddr:     mustGet("REDIS_ADDR"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       mustGetIntDefault("REDIS_DB", 0),

		RedisPoolSize:     mustGetIntDefault("REDIS_POOL_SIZE", 10),
		RedisMinIdleConns: mustGetIntDefault("REDIS_MIN_IDLE_CONNS", 3),
	}
}

func mustGet(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Missing required environment variable: %s", key)
	}
	return value
}

func mustGetIntDefault(key string, def int) int {
	str := os.Getenv(key)
	if str == "" {
		return def
	}
	i, err := strconv.Atoi(str)
	if err != nil {
		log.Fatalf("Invalid integer for %s: %s", key, str)
	}
	return i
}
