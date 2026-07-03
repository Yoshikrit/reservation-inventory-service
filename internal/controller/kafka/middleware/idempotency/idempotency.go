package idempotency

import (
	"context"
	"fmt"
	"time"

	"github.com/Yoshikrit/inventory/internal/controller/kafka/middleware"

	kafkago "github.com/segmentio/kafka-go"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

const ttl = 24 * time.Hour

func Idempotency(rdb *redis.Client) middleware.Middleware {
	return func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, msg kafkago.Message) error {
			key := fmt.Sprintf("kafka:processed:%s:%d:%d", msg.Topic, msg.Partition, msg.Offset)

			set, err := rdb.SetNX(ctx, key, 1, ttl).Result()
			if err != nil {
				log.Warn().Err(err).Str("key", key).Msg("kafka: idempotency check failed, processing anyway")
				return next(ctx, msg)
			}

			if !set {
				log.Info().
					Str("topic", msg.Topic).
					Int("partition", msg.Partition).
					Int64("offset", msg.Offset).
					Msg("kafka: skipping duplicate message")
				return nil
			}

			if err := next(ctx, msg); err != nil {
				rdb.Del(context.WithoutCancel(ctx), key)
				return err
			}
			return nil
		}
	}
}
