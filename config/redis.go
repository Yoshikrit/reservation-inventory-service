package config

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	RedisHost string `env:"REDIS_HOST,required"`
	RedisPort int    `env:"REDIS_PORT,required"`
}

func InitRedis(cfg RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
	})
}
