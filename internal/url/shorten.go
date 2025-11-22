package url

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	day              = 24 * time.Hour
	redisShortPrefix = "short:"
	redisLongPrefix  = "url:"
)

func ShortenURL(ctx context.Context, rdb *redis.Client, rawURL string) (string, error) {
	if err := validateURL(rawURL); err != nil {
		return "", err
	}

	normURL, err := normalizeURL(rawURL)
	if err != nil {
		return "", errors.New("failed to normalize URL")
	}

	longHash := hashURL(normURL)

	// Check if URL already exists
	existingCode, err := rdb.Get(ctx, redisLongPrefix+longHash).Result()
	if err == nil {
		rdb.Expire(ctx, redisLongPrefix+longHash, day)
		rdb.Expire(ctx, redisShortPrefix+existingCode, day)
		return existingCode, nil
	}
	if err != redis.Nil {
		return "", err
	}

	var code string
	for {
		code, err = generateCode()
		if err != nil {
			return "", err
		}
		ok, err := rdb.SetNX(ctx, redisShortPrefix+code, normURL, day).Result()
		if err != nil {
			return "", err
		}
		if ok {
			break
		}
	}

	if err := rdb.Set(ctx, redisLongPrefix+longHash, code, day).Err(); err != nil {
		return "", err
	}

	return code, nil
}

func GetURL(ctx context.Context, rdb *redis.Client, code string) (string, error) {
	return rdb.Get(ctx, redisShortPrefix+code).Result()
}

func generateCode() (string, error) {
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func hashURL(longURL string) string {
	sum := sha256.Sum256([]byte(longURL))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}
