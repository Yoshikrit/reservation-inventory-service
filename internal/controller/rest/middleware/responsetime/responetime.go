package responsetime

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/responsetime"
)

func ResponseTime() fiber.Handler {
	return responsetime.New()
}
