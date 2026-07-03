package rest

import (
	"errors"

	"github.com/Yoshikrit/inventory/internal/pkg/apperror"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

func ErrorHandler() fiber.ErrorHandler {
	return func(c fiber.Ctx, err error) error {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			return c.Status(categoryToHTTPCode(appErr.Category)).JSON(fiber.Map{
				"code":    appErr.CodeNumber,
				"message": appErr.Description,
			})
		}

		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			return c.Status(fiberErr.Code).JSON(fiber.Map{
				"code":    fiberErr.Code,
				"message": fiberErr.Message,
			})
		}

		log.Error().Err(err).Msg("internal server error")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "internal server error",
		})
	}
}

func categoryToHTTPCode(category apperror.Category) int {
	switch category {
	case apperror.CategoryBadRequest:
		return fiber.StatusBadRequest
	case apperror.CategoryUnauthorized:
		return fiber.StatusUnauthorized
	case apperror.CategoryNotFound:
		return fiber.StatusNotFound
	case apperror.CategoryConflict:
		return fiber.StatusConflict
	case apperror.CategoryUnprocessable:
		return fiber.StatusUnprocessableEntity
	default:
		return fiber.StatusInternalServerError
	}
}
