package inventory

import (
	inventorySrv "github.com/Yoshikrit/inventory/internal/service/inventory"

	"github.com/gofiber/fiber/v3"
)

func (c *InventoryController) GetProducts(ctx fiber.Ctx) error {
	products, err := c.inventoryService.GetProducts(ctx.Context())
	if err != nil {
		return err
	}

	response := parseProductListResponse(products)
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func parseProductListResponse(resp inventorySrv.GetProductsResponse) ProductListResponse {
	var productListResponse []ProductResponse
	for _, productEntity := range resp.Products {
		productResponse := ProductResponse{
			ProductID:   productEntity.ProductID,
			Name:        productEntity.Name,
			Description: productEntity.Description,
			Price:       productEntity.Price,
			Quantity:    productEntity.Quantity,
		}
		productListResponse = append(productListResponse, productResponse)
	}

	response := ProductListResponse{
		ProductList: productListResponse,
	}
	return response
}
