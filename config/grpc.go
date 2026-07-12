package config

import (
	"github.com/Yoshikrit/inventory/internal/controller/grpc/interceptor"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GrpcConfig struct {
	GrpcPort       int  `env:"GRPC_PORT,required"`
	GrpcReflection bool `env:"GRPC_REFLECTION" envDefault:"false"`
}

func InitGrpc(cfg GrpcConfig) *grpc.Server {
	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			interceptor.RecoveryUnary(),
			interceptor.TraceUnary(),
			interceptor.LoggingUnary(),
			interceptor.ValidationUnary(),
			interceptor.ErrorUnary(),
		),
	)
	if cfg.GrpcReflection {
		reflection.Register(grpcServer)
	}
	return grpcServer
}
