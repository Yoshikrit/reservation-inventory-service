package interceptor

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RecoveryUnary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				log.Error().
					Str("method", info.FullMethod).
					Str("panic", fmt.Sprintf("%v", r)).
					Bytes("stack", debug.Stack()).
					Msg("grpc: panic recovered")
				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()
		return handler(ctx, req)
	}
}
