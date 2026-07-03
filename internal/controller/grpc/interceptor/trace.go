package interceptor

import (
	"context"
	"fmt"

	"github.com/Yoshikrit/inventory/internal/entity"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TraceUnary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md, _ := metadata.FromIncomingContext(ctx)
		traceID := ""
		if vals := md.Get("x-request-id"); len(vals) > 0 {
			traceID = vals[0]
		}
		ctx = context.WithValue(ctx, entity.ContextKeyEvent, fmt.Sprintf("Entry=gRPC : %s", info.FullMethod))
		ctx = context.WithValue(ctx, entity.ContextKeyTraceID, traceID)
		return handler(ctx, req)
	}
}
