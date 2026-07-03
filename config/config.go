package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	RestConfig     RestConfig
	DatabaseConfig DatabaseConfig
	RedisConfig    RedisConfig
	GrpcConfig     GrpcConfig
	KafkaConfig    KafkaConfig
}

func Load() (Config, error) {
	cfg := Config{}
	_ = godotenv.Load()

	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
