package interceptor

import (
	"context"
	"time"

	"github.com/Yoshikrit/inventory/internal/entity"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func LoggingUnary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		code := codes.OK
		if err != nil {
			code = status.Code(err)
		}

		traceID, _ := ctx.Value(entity.ContextKeyTraceID).(string)
		event := log.Info().
			Str("method", info.FullMethod).
			Str("requestId", traceID).
			Str("code", code.String()).
			Dur("duration", duration)

		if err != nil {
			event = log.Error().
				Str("method", info.FullMethod).
				Str("requestId", traceID).
				Str("code", code.String()).
				Dur("duration", duration).
				Str("error", status.Convert(err).Message())
		}

		event.Msg("gRPC")

		return resp, err
	}
}
