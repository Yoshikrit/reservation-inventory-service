package middleware

import (
	"github.com/Yoshikrit/inventory/internal/controller/rest/middleware/cors"
	"github.com/Yoshikrit/inventory/internal/controller/rest/middleware/helmet"
	"github.com/Yoshikrit/inventory/internal/controller/rest/middleware/logger"
	"github.com/Yoshikrit/inventory/internal/controller/rest/middleware/recover"
	"github.com/Yoshikrit/inventory/internal/controller/rest/middleware/responsetime"
	"github.com/Yoshikrit/inventory/internal/controller/rest/middleware/trace"

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
