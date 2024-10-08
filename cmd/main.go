package main

import (
	"fmt"
	"metodichka/internal/cache"
	"metodichka/internal/config"
	"metodichka/internal/logger"
	"metodichka/internal/service"
	"metodichka/internal/transport"
)

func main() {
	cfg := config.Config{
		Transport: config.Transport{
			Host: "",
			Port: "8080",
		},
		Cache: config.Cache{
			Addr: "localhost:6379",
		},
	}

	c := cache.NewRedisCache(cfg.Cache)

	s := service.NewService(c)

	tr := transport.NewTransport(cfg.Transport, s)

	defer func() {
		if err := tr.Close(); err != nil {
			logger.GetLogger().Fatalln(err)
		}
	}()

	if err := tr.Run(); err != nil {
		logger.GetLogger().Println(fmt.Errorf("failed to run transport. %w", err))
	}
}
