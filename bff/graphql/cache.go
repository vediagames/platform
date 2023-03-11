package graphql

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
)

type Cache struct {
	client    redis.UniversalClient
	ttl       time.Duration
	apqPrefix string
}

func NewCache(ctx context.Context, redisAddress string, ttl time.Duration) Cache {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddress,
	})

	err := client.Ping(ctx).Err()
	if err != nil {
		panic(fmt.Errorf("failed to ping: %w", err))
	}

	return Cache{
		client: client,
		ttl:    ttl,
	}
}

func (c Cache) Add(ctx context.Context, key string, value interface{}) {
	key = fmt.Sprintf("%s:%s", c.apqPrefix, key)

	_, err := c.client.Set(ctx, key, value, c.ttl).Result()
	if err != nil {
		zerolog.Ctx(ctx).Err(fmt.Errorf("failed to set cache for key %q: %w", key, err))
	}
}

func (c Cache) Get(ctx context.Context, key string) (interface{}, bool) {
	key = fmt.Sprintf("%s:%s", c.apqPrefix, key)

	s, err := c.client.Get(ctx, key).Result()
	if err != nil {
		zerolog.Ctx(ctx).Err(fmt.Errorf("failed to get cache for key %q: %w", key, err))
		return nil, false
	}

	return s, true
}
