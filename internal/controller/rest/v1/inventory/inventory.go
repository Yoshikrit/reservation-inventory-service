package inventory

import (
	inventorySrv "inventory/internal/service/inventory"

	"github.com/gofiber/fiber/v3"
)

type InventoryController struct {
	inventoryService inventorySrv.InventoryService
}

func NewInventoryController(inventoryService inventorySrv.InventoryService) *InventoryController {
	return &InventoryController{inventoryService: inventoryService}
}

func (c *InventoryController) RegisterRoutes(router fiber.Router) {
	router.Post("/", c.CreateProduct)
	router.Get("/", c.GetProducts)
	router.Get("/:product_id", c.GetProductByID)
}
