package trace

import (
	"context"
	"fmt"

	"inventory/internal/controller/kafka/middleware"
	"inventory/internal/entity"

	"github.com/google/uuid"
	kafka "github.com/segmentio/kafka-go"
)

func Trace() middleware.Middleware {
	return func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, msg kafka.Message) error {
			traceID := ""
			for _, h := range msg.Headers {
				if h.Key == "X-Request-ID" {
					traceID = string(h.Value)
					break
				}
			}
			if traceID == "" {
				traceID = uuid.New().String()
			}
			ctx = context.WithValue(ctx, entity.ContextKeyTraceID, traceID)
			ctx = context.WithValue(ctx, entity.ContextKeyEvent, fmt.Sprintf("Entry=Kafka : %s", msg.Topic))
			return next(ctx, msg)
		}
	}
}
