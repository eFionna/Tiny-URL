package url

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func setupRedis(t *testing.T) (*redis.Client, func()) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	return rdb, func() { mr.Close() }
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		raw   string
		valid bool
	}{
		{"https://example.com", true},
		{"http://example.com", true},
		{"ftp://example.com", false},
		{"javascript:alert(1)", false},
		{"", false},
	}

	for _, tt := range tests {
		err := validateURL(tt.raw)
		if tt.valid {
			assert.NoError(t, err, tt.raw)
		} else {
			assert.Error(t, err, tt.raw)
		}
	}
}

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{"https://Example.com/", "https://example.com/"},
		{"https://example.com/path/", "https://example.com/path"},
		{"https://example.com/path#fragment", "https://example.com/path"},
	}

	for _, tt := range tests {
		got, err := normalizeURL(tt.input)
		assert.NoError(t, err)
		assert.Equal(t, tt.output, got)
	}
}

func TestShortenAndGetURL(t *testing.T) {
	rdb, cleanup := setupRedis(t)
	defer cleanup()
	ctx := context.Background()

	raw := "https://example.com/test"
	code, err := ShortenURL(ctx, rdb, raw)
	assert.NoError(t, err)
	assert.NotEmpty(t, code)

	long, err := GetURL(ctx, rdb, code)
	assert.NoError(t, err)
	assert.Contains(t, long, "example.com/test")

	code2, err := ShortenURL(ctx, rdb, raw)
	assert.NoError(t, err)
	assert.Equal(t, code, code2)
}

func TestShortenURLDefensive(t *testing.T) {
	rdb, cleanup := setupRedis(t)
	defer cleanup()
	ctx := context.Background()

	tests := []struct {
		rawURL      string
		shouldError bool
	}{
		{"", true},                        // empty URL
		{"ftp://example.com", true},       // invalid scheme
		{"javascript:alert(1)", true},     // unsafe scheme
		{"data:text/plain,hello", true},   // unsafe scheme
		{"vbscript:msgbox('hi')", true},   // unsafe scheme
		{"http://valid.com", false},       // valid
		{"https://valid.com/path", false}, // valid
	}

	for _, tt := range tests {
		code, err := ShortenURL(ctx, rdb, tt.rawURL)
		if tt.shouldError {
			assert.Error(t, err, tt.rawURL)
			assert.Empty(t, code)
		} else {
			assert.NoError(t, err, tt.rawURL)
			assert.NotEmpty(t, code)
		}
	}
}
