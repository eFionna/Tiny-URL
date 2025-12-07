package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	os.Setenv("PORT", "8080")
	os.Setenv("BASE_URL", "http://localhost:8080")
	os.Setenv("REDIS_ADDR", "localhost:6379")

	cfg := Load()
	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "http://localhost:8080", cfg.BaseURL)
	assert.Equal(t, "localhost:6379", cfg.RedisAddr)
}
