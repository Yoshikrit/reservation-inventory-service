package logger

import (
	"context"
	"time"

	"inventory/internal/controller/kafka/middleware"
	"inventory/internal/entity"

	kafka "github.com/segmentio/kafka-go"
	"github.com/rs/zerolog/log"
)

func Logger() middleware.Middleware {
	return func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, msg kafka.Message) error {
			start := time.Now()
			err := next(ctx, msg)
			duration := time.Since(start)
			traceID, _ := ctx.Value(entity.ContextKeyTraceID).(string)

			if err != nil {
				log.Error().
					Str("topic", msg.Topic).
					Int64("offset", msg.Offset).
					Int("partition", msg.Partition).
					Str("requestId", traceID).
					Dur("duration", duration).
					Str("error", err.Error()).
					Msg("Kafka")
				return err
			}

			log.Info().
				Str("topic", msg.Topic).
				Int64("offset", msg.Offset).
				Int("partition", msg.Partition).
				Str("requestId", traceID).
				Dur("duration", duration).
				Msg("Kafka")
			return nil
		}
	}
}
