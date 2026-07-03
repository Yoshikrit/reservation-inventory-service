package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"

	"github.com/Yoshikrit/inventory/config"
	grpcv1 "github.com/Yoshikrit/inventory/internal/controller/grpc/v1"
	"github.com/Yoshikrit/inventory/internal/pkg/logger"
)

func main() {
	logger.Init()
	log.Info().Msg("inventory-grpc: starting")

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("inventory-grpc: failed to load config")
	}

	db, err := config.InitDatabase(cfg.DatabaseConfig.DatabaseUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("inventory-grpc: failed to connect database")
	}
	if err := config.MigrateDatabase(db); err != nil {
		log.Fatal().Err(err).Msg("inventory-grpc: failed to migrate database")
	}

	redis := config.InitRedis(cfg.RedisConfig)

	grpcServer := config.InitGrpc(cfg.GrpcConfig)
	grpcv1.NewGRPCRouter(grpcServer, db, redis)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GrpcConfig.GrpcPort))
	if err != nil {
		log.Fatal().Err(err).Msg("inventory-grpc: failed to listen")
	}

	go func() {
		log.Info().Int("port", cfg.GrpcConfig.GrpcPort).Msg("inventory-grpc: listening")
		if err := grpcServer.Serve(lis); err != nil {
			log.Error().Err(err).Msg("inventory-grpc: server stopped")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("inventory-grpc: shutting down")
	grpcServer.GracefulStop()
}
