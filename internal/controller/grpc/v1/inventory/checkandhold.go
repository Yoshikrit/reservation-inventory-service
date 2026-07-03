package inventory

import (
	"context"

	"github.com/Yoshikrit/inventory/internal/controller/grpc/v1/inventory/pb"
	inventorySrv "github.com/Yoshikrit/inventory/internal/service/inventory"
)

func (c *InventoryController) CheckAndHold(ctx context.Context, request *pb.CheckAndHoldRequest) (*pb.CheckAndHoldResponse, error) {
	result, err := c.inventoryService.CheckAndHold(ctx, &inventorySrv.CheckAndHoldRequest{
		ProductID: request.ProductId,
		Quantity:  uint(request.Quantity),
	})
	if err != nil {
		return nil, err
	}
	return &pb.CheckAndHoldResponse{
		Available: result.Available,
		Product: &pb.ProductResponse{
			ProductId:   result.ProductID,
			Name:        result.Name,
			Description: result.Description,
			Price:       result.Price,
			Quantity:    uint64(result.Quantity),
		},
	}, nil
}
