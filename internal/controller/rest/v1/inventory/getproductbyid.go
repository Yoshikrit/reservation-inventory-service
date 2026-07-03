package inventory

import (
	inventorySrv "inventory/internal/service/inventory"
	"inventory/internal/pkg/apperror"

	"github.com/gofiber/fiber/v3"
)

func (c *InventoryController) GetProductByID(ctx fiber.Ctx) error {
	var param GetProductByIDRequest
	if err := ctx.Bind().URI(&param); err != nil {
		return apperror.NewError(40000000, err)
	}

	product, err := c.inventoryService.GetProductByID(ctx.Context(), &inventorySrv.GetProductRequest{
		ProductID: param.ProductID,
	})
	if err != nil {
		return err
	}

	response := parseProductResponse(product)
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func parseProductResponse(resp inventorySrv.GetProductResponse) ProductResponse {
	return ProductResponse{
		ProductID:   resp.ProductID,
		Name:        resp.Name,
		Description: resp.Description,
		Price:       resp.Price,
		Quantity:    resp.Quantity,
	}
}
