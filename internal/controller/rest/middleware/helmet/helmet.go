package helmet

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/helmet"
)

func Helmet() fiber.Handler {
	return helmet.New()
}