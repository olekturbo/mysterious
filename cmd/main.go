// @title Mysterious API
// @version 1.0
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/olekturbo/mysterious/docs"
	"github.com/olekturbo/mysterious/internal/api"
	"github.com/olekturbo/mysterious/internal/cache"
	"github.com/olekturbo/mysterious/internal/predictor"
	"github.com/olekturbo/mysterious/internal/service"
	"github.com/olekturbo/mysterious/internal/storage"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Config struct {
	Port  string
	Redis struct {
		Addr     string
		Password string
	}
	Limit struct {
		Count int64
		Exp   time.Duration
	}
	MySQL struct {
		DSN string
	}
	Token struct {
		Key string
		Exp time.Duration
	}
	Predict struct {
		URL string
	}
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		panic(err)
	}

	db, err := storage.NewMySQL(cfg.MySQL.DSN)
	if err != nil {
		panic(err)
	}

	handler := api.NewHandler(
		service.NewID(),
		service.NewCache(cache.NewRedis(cfg.Redis.Addr, cfg.Redis.Password),
			cfg.Limit.Count,
			cfg.Limit.Exp,
		),
		service.NewUser(
			db,
			service.NewID(),
		),
		service.NewToken(cfg.Token.Key, cfg.Token.Exp),
		service.NewPredict(predictor.NewHTTP(http.DefaultClient, cfg.Predict.URL)),
		api.NewCookieManager(),
	)

	router := chi.NewRouter()
	router.Get("/", handler.Home)
	router.Post("/register", handler.Register)
	router.Post("/login", handler.Login)
	router.Get("/swagger/*", httpSwagger.WrapHandler)

	router.Route("/private", func(r chi.Router) {
		tokenAuth := jwtauth.New("HS256", []byte(cfg.Token.Key), nil)
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))

		r.Get("/", handler.PrivateHome)
	})

	router.Post("/predict", handler.Predict)

	err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), router)
	if err != nil {
		panic(err)
	}
}
