package config

import (
	"github.com/Yoshikrit/inventory/internal/pkg/validator"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
)

type RestConfig struct {
	RestPort string `env:"REST_PORT,required"`
}

func NewRestConfig(errorHandler fiber.ErrorHandler) fiber.Config {
	return fiber.Config{
		JSONEncoder:     sonic.Marshal,
		JSONDecoder:     sonic.Unmarshal,
		StructValidator: validator.New(),
		ErrorHandler:    errorHandler,
	}
}
