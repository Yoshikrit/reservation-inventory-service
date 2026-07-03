package inventory

import (
	"github.com/Yoshikrit/inventory/internal/controller/grpc/v1/inventory/pb"
	inventoryService "github.com/Yoshikrit/inventory/internal/service/inventory"
)

type InventoryController struct {
	pb.UnimplementedInventoryServiceServer
	inventoryService inventoryService.InventoryService
}

func NewInventoryController(inventoryService inventoryService.InventoryService) *InventoryController {
	return &InventoryController{inventoryService: inventoryService}
}
