package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kelseyhightower/envconfig"
	"github.com/olekturbo/mysterious/internal/api"
	"github.com/olekturbo/mysterious/internal/cache"
	"github.com/olekturbo/mysterious/internal/service"
)

type Config struct {
	Port  string
	Redis Redis
	Limit int64
	Exp   time.Duration
}

type Redis struct {
	Addr     string
	Password string
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		panic(err)
	}

	handler := api.NewHandler(
		service.NewID(),
		service.NewCache(cache.NewRedis(cfg.Redis.Addr, cfg.Redis.Password),
			cfg.Limit,
			cfg.Exp),
		api.NewCookieManager(),
	)

	router := chi.NewRouter()
	router.Get("/", handler.Home)

	err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), router)
	if err != nil {
		panic(err)
	}
}
