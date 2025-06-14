package api

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrLimitReached = errors.New("limit reached")

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string, exp time.Duration) error
	Incr(ctx context.Context, key string) (int64, error)
}

type Service struct {
	cache Cache
	limit int64
	exp   time.Duration
}

func NewService(cache Cache, limit int64, exp time.Duration) *Service {
	return &Service{
		cache: cache,
		limit: limit,
		exp:   exp,
	}
}

func (s *Service) Limit(ctx context.Context, key string) error {
	val, err := s.cache.Get(ctx, key)
	if err != nil {
		return err
	}

	if val == "" {
		err = s.cache.Set(ctx, key, "1", s.exp)
		if err != nil {
			return err
		}
		return nil
	}

	rate, err := s.cache.Incr(ctx, key)
	if err != nil {
		return err
	}

	if rate > s.limit {
		return ErrLimitReached
	}

	return nil
}

func (s *Service) GenerateID() string {
	return uuid.NewString()
}
