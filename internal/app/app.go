package app

import (
	"context"
	"html/template"
	"time"

	"github.com/eFionna/Tiny-URL/internal/config"
	"github.com/redis/go-redis/v9"
)

//var tmpl = template.Must(template.ParseFiles("web/index.html"))

type App struct {
	RDB    *redis.Client
	Config *config.Config
	Ctx    context.Context
	Tmpl   *template.Template
}

func NewApp(cfg *config.Config, tmpl *template.Template) (*App, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:            cfg.RedisAddr,
		Password:        cfg.RedisPassword,
		DB:              cfg.RedisDB,
		PoolSize:        cfg.RedisPoolSize,
		MinIdleConns:    cfg.RedisMinIdleConns,
		ConnMaxIdleTime: 5 * time.Minute,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &App{
		RDB:    rdb,
		Config: cfg,
		Ctx:    ctx,
		Tmpl:   tmpl,
	}, nil
}
