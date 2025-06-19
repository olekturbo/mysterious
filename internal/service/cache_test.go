package service

import (
	"context"
	"errors"
	"testing"
	"time"
)

type mockCache struct {
	getFn  func(ctx context.Context, key string) (string, error)
	setFn  func(ctx context.Context, key string, value string, exp time.Duration) error
	incrFn func(ctx context.Context, key string) (int64, error)
}

func (m *mockCache) Get(ctx context.Context, key string) (string, error) {
	return m.getFn(ctx, key)
}

func (m *mockCache) Set(ctx context.Context, key string, value string, exp time.Duration) error {
	return m.setFn(ctx, key, value, exp)
}

func (m *mockCache) Incr(ctx context.Context, key string) (int64, error) {
	return m.incrFn(ctx, key)
}

func TestCache_Limit(t *testing.T) {
	tests := []struct {
		name      string
		mockCache mockCache
		key       string
		wantErr   bool
	}{
		{
			name: "new key sets correctly",
			mockCache: mockCache{
				getFn: func(ctx context.Context, key string) (string, error) {
					return "", nil
				},
				setFn: func(ctx context.Context, key string, value string, exp time.Duration) error {
					return nil
				},
			},
			key:     "test-key",
			wantErr: false,
		},
		{
			name: "increment within limit",
			mockCache: mockCache{
				getFn: func(ctx context.Context, key string) (string, error) {
					return "1", nil
				},
				incrFn: func(ctx context.Context, key string) (int64, error) {
					return 2, nil
				},
			},
			key:     "test-key",
			wantErr: false,
		},
		{
			name: "increment exceeds limit",
			mockCache: mockCache{
				getFn: func(ctx context.Context, key string) (string, error) {
					return "2", nil
				},
				incrFn: func(ctx context.Context, key string) (int64, error) {
					return 6, nil
				},
			},
			key:     "test-key",
			wantErr: true,
		},
		{
			name: "get error",
			mockCache: mockCache{
				getFn: func(ctx context.Context, key string) (string, error) {
					return "", errors.New("get error")
				},
			},
			key:     "test-key",
			wantErr: true,
		},
		{
			name: "set error on new key",
			mockCache: mockCache{
				getFn: func(ctx context.Context, key string) (string, error) {
					return "", nil
				},
				setFn: func(ctx context.Context, key string, value string, exp time.Duration) error {
					return errors.New("set error")
				},
			},
			key:     "test-key",
			wantErr: true,
		},
		{
			name: "increment error",
			mockCache: mockCache{
				getFn: func(ctx context.Context, key string) (string, error) {
					return "1", nil
				},
				incrFn: func(ctx context.Context, key string) (int64, error) {
					return 0, errors.New("increment error")
				},
			},
			key:     "test-key",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := &Cache{
				cache: &tt.mockCache,
				limit: 5,
				exp:   time.Second,
			}

			err := cache.Limit(context.Background(), tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected wantErr %v, got err %v", tt.wantErr, err)
			}
		})
	}
}
