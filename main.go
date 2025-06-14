package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	"github.com/olekturbo/mysterious/internal/cache"
)

type Config struct {
	Port  string
	Redis Redis
}

type Redis struct {
	Addr     string
	Password string
}

type Handler struct {
	s string
}

func (h *Handler) Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, h.s)
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		panic(err)
	}

	id := uuid.NewString()

	redisCache := cache.NewRedis(cfg.Redis.Addr, cfg.Redis.Password)
	err = redisCache.Set(context.TODO(), id, "Hello!")
	if err != nil {
		panic(err)
	}

	received, err := redisCache.Get(context.Background(), id)
	if err != nil {
		panic(err)
	}

	handler := Handler{
		s: received,
	}

	http.HandleFunc("/", handler.Hello)

	port := fmt.Sprintf(":%s", cfg.Port)
	fmt.Printf("Server starting at %s", port)
	http.ListenAndServe(port, nil)
}
