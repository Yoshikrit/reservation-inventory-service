package grpcutil

import (
	"context"

	"github.com/Yoshikrit/inventory/internal/entity"

	"google.golang.org/grpc/metadata"
)

// AppendTrace injects the trace ID from context into outgoing gRPC metadata.
// Call this before every gRPC client request so the downstream service receives the trace.
func AppendTrace(ctx context.Context) context.Context {
	traceID, _ := ctx.Value(entity.ContextKeyTraceID).(string)
	if traceID == "" {
		return ctx
	}
	return metadata.AppendToOutgoingContext(ctx, "x-request-id", traceID)
}
