package recovery

import (
	"context"
	"fmt"
	"runtime/debug"

	"inventory/internal/controller/kafka/middleware"

	kafka "github.com/segmentio/kafka-go"
	"github.com/rs/zerolog/log"
)

func Recovery() middleware.Middleware {
	return func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, msg kafka.Message) (err error) {
			defer func() {
				if r := recover(); r != nil {
					log.Error().
						Str("topic", msg.Topic).
						Int64("offset", msg.Offset).
						Str("panic", fmt.Sprintf("%v", r)).
						Str("stack", string(debug.Stack())).
						Msg("kafka: panic recovered")
					err = fmt.Errorf("panic: %v", r)
				}
			}()
			return next(ctx, msg)
		}
	}
}
