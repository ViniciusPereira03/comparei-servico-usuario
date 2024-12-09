package db

import (
	"comparei-servico-usuario/config"
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

func ConnectRedis(cfg *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		DB:   0,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao Redis: %v", err)
	}

	log.Println("Conexão com Redis estabelecida com sucesso")
	return rdb, nil
}
