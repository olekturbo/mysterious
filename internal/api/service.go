package api

import (
	"context"

	"github.com/google/uuid"
)

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string) error
}

type Service struct {
	cache Cache
}

func NewService(cache Cache) *Service {
	return &Service{
		cache: cache,
	}
}

func (s *Service) Get(ctx context.Context, key string) (string, error) {
	return s.cache.Get(ctx, key)
}

func (s *Service) Set(ctx context.Context, key, value string) error {
	return s.cache.Set(ctx, key, value)
}

func (s *Service) GenerateID() string {
	return uuid.NewString()
}
