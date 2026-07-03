package kafka

import (
	"context"
	"sync"

	"inventory/config"
	"inventory/internal/controller/kafka/middleware"
	kafkalogger "inventory/internal/controller/kafka/middleware/logger"
	kafkarecovery "inventory/internal/controller/kafka/middleware/recovery"
	kafkatrace "inventory/internal/controller/kafka/middleware/trace"

	"github.com/rs/zerolog/log"
)

func Start(ctx context.Context, cfg config.KafkaConfig, consumers []consumerConfig) {
	var wg sync.WaitGroup
	for _, c := range consumers {
		wg.Add(1)
		go func(c consumerConfig) {
			defer wg.Done()
			startConsumer(ctx, cfg, c)
		}(c)
	}
	wg.Wait()
}

func startConsumer(ctx context.Context, cfg config.KafkaConfig, c consumerConfig) {
	reader := config.NewKafkaReader(cfg, c.topic, c.groupID)
	defer reader.Close()

	handler := middleware.Chain(
		c.handler,
		kafkarecovery.Recovery(),
		kafkatrace.Trace(),
		kafkalogger.Logger(),
	)

	log.Info().Str("topic", c.topic).Str("group", c.groupID).Msg("kafka: consumer started")

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Info().Str("topic", c.topic).Msg("kafka: consumer stopped")
				return
			}
			log.Error().Err(err).Str("topic", c.topic).Msg("kafka: failed to read message")
			continue
		}

		if err := handler(ctx, msg); err != nil {
			log.Error().Err(err).Str("topic", msg.Topic).Int64("offset", msg.Offset).Msg("kafka: failed to handle message")
		}
	}
}
