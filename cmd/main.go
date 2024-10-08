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
	trCfg := config.Transport{
		Host: "",
		Port: "8080",
	}

	cacheCfg := config.Cache{
		Addr: "localhost:6379",
	}
	c := cache.NewRedisCache(cacheCfg)

	s := service.NewService(c)

	tr := transport.NewTransport(trCfg, s)

	defer func() {
		if err := tr.Close(); err != nil {
			logger.GetLogger().Fatalln(err)
		}
	}()

	if err := tr.Run(); err != nil {
		logger.GetLogger().Println(fmt.Errorf("failed to run transport. %w", err))
	}
}
