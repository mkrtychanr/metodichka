package cache

import (
	"context"
	"fmt"
	"metodichka/internal/config"

	"github.com/go-redis/redis/v8"
)

const redisExists = 1
const timeWithNoTTLInRedis = 0

type redisCache struct {
	client *redis.Client
}

func NewRedisCache(cfg config.Cache) Cache {
	cl := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       0,
	})

	return &redisCache{
		client: cl,
	}
}

func (c *redisCache) Get(ctx context.Context, key string) (string, bool, error) {
	exists, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return "", false, fmt.Errorf("failed to check value in redis. %w", err)
	}

	if exists != redisExists {
		return "", false, nil
	}

	value, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return "", false, fmt.Errorf("failed to get value from redis. %w", err)
	}

	return value, true, nil
}

func (c *redisCache) Set(ctx context.Context, key, value string) error {
	if err := c.client.Set(ctx, key, value, timeWithNoTTLInRedis).Err(); err != nil {
		return fmt.Errorf("failed to set value in redis. %w", err)
	}

	return nil
}

func (c *redisCache) Close() error {
	c.client.Close()

	return nil
}
