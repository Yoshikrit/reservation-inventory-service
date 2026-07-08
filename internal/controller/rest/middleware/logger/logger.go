package logger

import (
	"github.com/gofiber/contrib/v3/zerolog"
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

var skipPaths = map[string]bool{
	"/metrics": true,
	"/livez":   true,
	"/readyz":  true,
}

func Logger() fiber.Handler {
	return zerolog.New(zerolog.Config{
		Logger: &log.Logger,
		Fields: []string{"latency", "status", "method", "url", "ip", "requestId", "error"},
		Skip:   func(c fiber.Ctx) bool { return skipPaths[c.Path()] },
	})
}
