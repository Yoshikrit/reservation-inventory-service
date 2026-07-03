package inventory

import (
	"context"

	"github.com/Yoshikrit/inventory/internal/controller/grpc/v1/inventory/pb"
	inventorySrv "github.com/Yoshikrit/inventory/internal/service/inventory"
)

func (c *InventoryController) GetProductByID(ctx context.Context, request *pb.GetProductByIDRequest) (*pb.ProductResponse, error) {
	product, err := c.inventoryService.GetProductByID(ctx, &inventorySrv.GetProductRequest{
		ProductID: request.ProductId,
	})
	if err != nil {
		return nil, err
	}
	return &pb.ProductResponse{
		ProductId:   product.ProductID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Quantity:    uint64(product.Quantity),
	}, nil
}
