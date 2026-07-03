package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"

	"github.com/Yoshikrit/inventory/config"
	kafka "github.com/Yoshikrit/inventory/internal/controller/kafka"
	"github.com/Yoshikrit/inventory/internal/pkg/logger"
)

func main() {
	logger.Init()
	log.Info().Msg("inventory-consumer: starting")

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("inventory-consumer: failed to load config")
	}

	db, err := config.InitDatabase(cfg.DatabaseConfig.DatabaseUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("inventory-consumer: failed to connect database")
	}
	if err := config.MigrateDatabase(db); err != nil {
		log.Fatal().Err(err).Msg("inventory-consumer: failed to migrate database")
	}

	redis := config.InitRedis(cfg.RedisConfig)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	consumers := kafka.NewKafkaRouter(cfg.KafkaConfig, db, redis)
	kafka.Start(ctx, cfg.KafkaConfig, consumers)

	log.Info().Msg("inventory-consumer: stopped")
}
