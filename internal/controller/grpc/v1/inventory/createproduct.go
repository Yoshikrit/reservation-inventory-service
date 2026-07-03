package inventory

import (
	"context"

	"github.com/Yoshikrit/inventory/internal/controller/grpc/v1/inventory/pb"
	inventorySrv "github.com/Yoshikrit/inventory/internal/service/inventory"
)

func (c *InventoryController) CreateProduct(ctx context.Context, request *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	srvReq := parseCreateRequestToService(request)
	if err := c.inventoryService.CreateProduct(ctx, &srvReq); err != nil {
		return nil, err
	}
	return &pb.CreateProductResponse{}, nil
}

func parseCreateRequestToService(req *pb.CreateProductRequest) inventorySrv.CreateProductRequest {
	return inventorySrv.CreateProductRequest{
		ProductID:   req.ProductId,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Quantity:    uint(req.Quantity),
	}
}
