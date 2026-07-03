package inventory

import inventorySrv "github.com/Yoshikrit/inventory/internal/service/inventory"

type KafkaInventoryController struct {
	inventoryService inventorySrv.InventoryService
}

func NewKafkaInventoryController(svc inventorySrv.InventoryService) *KafkaInventoryController {
	return &KafkaInventoryController{inventoryService: svc}
}
