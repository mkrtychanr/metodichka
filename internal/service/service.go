package service

import (
	"context"
	"fmt"
	"metodichka/internal/cache"
	"metodichka/internal/logger"
	"strconv"

	"github.com/shopspring/decimal"
)

type Service interface {
	GetFibonacci(int) string
	GetFibonacciCached(context.Context, int) string
	Close() error
}

type service struct {
	cache cache.Cache
}

func NewService(c cache.Cache) Service {
	return &service{
		cache: c,
	}
}

func (s *service) GetFibonacci(n int) string {
	if n == 0 {
		return "0"
	}

	if n == 1 {
		return "1"
	}

	first := decimal.NewFromInt(0)
	second := decimal.NewFromInt(1)

	for i := 2; i < n; i++ {
		sum := first.Add(second)

		if i%2 == 0 {
			first = sum
		} else {
			second = sum
		}
	}

	return first.Add(second).String()

}

func (s *service) GetFibonacciCached(ctx context.Context, n int) string {
	key := strconv.Itoa(n)
	value, ok, err := s.cache.Get(ctx, key)
	if err != nil {
		logger.GetLogger().Println(fmt.Errorf("failed to get cache. %w", err))

		return s.GetFibonacci(n)
	}

	if !ok {
		value = s.GetFibonacci(n)

		if err := s.cache.Set(ctx, key, value); err != nil {
			logger.GetLogger().Println(fmt.Errorf("failed to set value to cache. %w", err))
		}
	}

	return value
}

func (s *service) Close() error {
	if err := s.cache.Close(); err != nil {
		return fmt.Errorf("failed to close cache. %w", err)
	}

	return nil
}
