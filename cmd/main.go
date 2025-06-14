package main

import (
	"fmt"
	"net/http"

	"github.com/kelseyhightower/envconfig"
	"github.com/olekturbo/mysterious/internal/api"
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

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		panic(err)
	}

	handler := api.NewHandler(
		api.NewService(
			cache.NewRedis(cfg.Redis.Addr, cfg.Redis.Password),
		),
	)

	http.HandleFunc("/", handler.Home)

	port := fmt.Sprintf(":%s", cfg.Port)
	fmt.Printf("Server starting at %s", port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
}
