package logger

import (
	"context"
	"time"

	"github.com/Yoshikrit/inventory/internal/controller/kafka/middleware"
	"github.com/Yoshikrit/inventory/internal/entity"

	kafka "github.com/segmentio/kafka-go"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/trace"
)

func Logger() middleware.Middleware {
	return func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, msg kafka.Message) error {
			start := time.Now()
			err := next(ctx, msg)
			duration := time.Since(start)
			requestID, _ := ctx.Value(entity.ContextKeyTraceID).(string)

			e := log.Info()
			if err != nil {
				e = log.Error().Err(err)
			}

			e = e.
				Str("topic", msg.Topic).
				Int64("offset", msg.Offset).
				Int("partition", msg.Partition).
				Str("requestId", requestID).
				Dur("duration", duration)

			span := trace.SpanFromContext(ctx)
			if span.SpanContext().IsValid() {
				e = e.
					Str("traceId", span.SpanContext().TraceID().String()).
					Str("spanId", span.SpanContext().SpanID().String())
			}

			e.Msg("Kafka")
			return err
		}
	}
}
