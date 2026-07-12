package otel

import (
	"context"

	"github.com/Yoshikrit/inventory/internal/controller/kafka/middleware"

	kafka "github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("kafka.consumer")

func OtelTrace() middleware.Middleware {
	return func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, msg kafka.Message) error {
			carrier := propagation.MapCarrier{}
			for _, h := range msg.Headers {
				carrier[h.Key] = string(h.Value)
			}
			ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

			ctx, span := tracer.Start(ctx, "kafka.consume "+msg.Topic,
				oteltrace.WithSpanKind(oteltrace.SpanKindConsumer),
				oteltrace.WithAttributes(
					attribute.String("messaging.system", "kafka"),
					attribute.String("messaging.destination", msg.Topic),
					attribute.Int64("messaging.kafka.offset", msg.Offset),
				),
			)
			defer span.End()

			return next(ctx, msg)
		}
	}
}
