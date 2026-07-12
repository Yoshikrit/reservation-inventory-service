package logger

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/trace"
)

var skipPaths = map[string]bool{
	"/metrics": true,
	"/livez":   true,
	"/readyz":  true,
}

func Logger() fiber.Handler {
	return func(c fiber.Ctx) error {
		if skipPaths[c.Path()] {
			return c.Next()
		}

		start := time.Now()
		err := c.Next()
		statusCode := c.Response().StatusCode()

		event := log.Logger.Info()
		if statusCode >= 500 {
			event = log.Logger.Error()
		} else if statusCode >= 400 {
			event = log.Logger.Warn()
		}

		e := event.
			Dur("latency", time.Since(start)).
			Int("status", statusCode).
			Str("method", c.Method()).
			Str("url", c.Path()).
			Str("ip", c.IP()).
			Str("requestId", c.GetRespHeader("X-Request-ID"))

		span := trace.SpanFromContext(c.Context())
		if span.SpanContext().IsValid() {
			e = e.
				Str("traceId", span.SpanContext().TraceID().String()).
				Str("spanId", span.SpanContext().SpanID().String())
		}

		if err != nil {
			e = e.Err(err)
		}

		e.Msg("Success")

		return err
	}
}
