package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"

	"github.com/Yoshikrit/inventory/config"
	"github.com/Yoshikrit/inventory/internal/controller/rest"
	"github.com/Yoshikrit/inventory/internal/pkg/logger"
)

func main() {
	logger.Init()
	log.Info().Msg("inventory-api: starting")

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("inventory-api: failed to load config")
	}

	db, err := config.InitDatabase(cfg.DatabaseConfig.DatabaseUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("inventory-api: failed to connect database")
	}

	redis := config.InitRedis(cfg.RedisConfig)

	if err := config.MigrateDatabase(db); err != nil {
		log.Fatal().Err(err).Msg("inventory-api: failed to migrate database")
	}

	app := fiber.New(config.NewRestConfig(rest.ErrorHandler()))
	rest.NewRestRouter(app, db, redis)

	go func() {
		if err := app.Listen(":" + cfg.RestConfig.RestPort); err != nil {
			log.Error().Err(err).Msg("inventory-api: server stopped")
		}
	}()

	log.Info().Str("port", cfg.RestConfig.RestPort).Msg("inventory-api: listening")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("inventory-api: shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = app.ShutdownWithContext(ctx)
}
