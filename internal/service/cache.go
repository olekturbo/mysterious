package service

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var ErrLimitReached = errors.New("limit reached")

type cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string, exp time.Duration) error
	Incr(ctx context.Context, key string) (int64, error)
}

type Cache struct {
	cache cache
	limit int64
	exp   time.Duration
}

func NewCache(cache cache, limit int64, exp time.Duration) *Cache {
	return &Cache{
		cache: cache,
		limit: limit,
		exp:   exp,
	}
}

func (c *Cache) Limit(ctx context.Context, key string) error {
	key = fmt.Sprintf("limit:%s", key)

	val, err := c.cache.Get(ctx, key)
	if err != nil {
		return err
	}

	if val == "" {
		err = c.cache.Set(ctx, key, "1", c.exp)
		if err != nil {
			return err
		}
		return nil
	}

	rate, err := c.cache.Incr(ctx, key)
	if err != nil {
		return err
	}

	if rate > c.limit {
		return ErrLimitReached
	}

	return nil
}
