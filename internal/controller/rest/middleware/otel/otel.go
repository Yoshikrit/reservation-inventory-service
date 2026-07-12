package otel

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("fiber.http")

func Trace() fiber.Handler {
	return func(c fiber.Ctx) error {
		propagator := otel.GetTextMapPropagator()
		ctx := propagator.Extract(c.Context(), propagation.HeaderCarrier(c.GetReqHeaders()))

		spanName := c.Method() + " " + c.Route().Path
		ctx, span := tracer.Start(ctx, spanName,
			oteltrace.WithSpanKind(oteltrace.SpanKindServer),
			oteltrace.WithAttributes(
				attribute.String("http.method", c.Method()),
				attribute.String("http.url", c.Path()),
				attribute.String("http.host", c.Hostname()),
			),
		)
		defer span.End()

		c.SetContext(ctx)

		err := c.Next()

		statusCode := c.Response().StatusCode()
		span.SetAttributes(attribute.Int("http.status_code", statusCode))
		if statusCode >= 500 {
			span.SetStatus(codes.Error, strconv.Itoa(statusCode))
		}

		return err
	}
}
