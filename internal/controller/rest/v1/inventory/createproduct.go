package inventory

import (
	"inventory/internal/pkg/apperror"

	inventorySrv "inventory/internal/service/inventory"

	"github.com/gofiber/fiber/v3"
)

func (c *InventoryController) CreateProduct(ctx fiber.Ctx) error {
	var request CreateProductRequest
	if err := ctx.Bind().JSON(&request); err != nil {
		return apperror.NewError(40000000, err)
	}

	requestToSrv := parseCreateRequestToService(request)
	if err := c.inventoryService.CreateProduct(ctx.Context(), &requestToSrv); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{})
}

func parseCreateRequestToService(request CreateProductRequest) inventorySrv.CreateProductRequest {
	return inventorySrv.CreateProductRequest{
		ProductID:   request.ProductID,
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
		Quantity:    request.Quantity,
	}
}
