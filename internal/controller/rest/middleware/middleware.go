package middleware

import (
	"inventory/internal/controller/rest/middleware/cors"
	"inventory/internal/controller/rest/middleware/helmet"
	"inventory/internal/controller/rest/middleware/logger"
	"inventory/internal/controller/rest/middleware/recover"
	"inventory/internal/controller/rest/middleware/responsetime"
	"inventory/internal/controller/rest/middleware/trace"

	"github.com/gofiber/fiber/v3"
)

func NewMiddleware() []fiber.Handler {
	return []fiber.Handler{
		recover.Recover(),
		logger.Logger(),
		cors.Cors(),
		responsetime.ResponseTime(),
		helmet.Helmet(),
		trace.Trace(),
	}
}
