package inventory

import (
	"context"

	"github.com/Yoshikrit/inventory/internal/controller/grpc/v1/inventory/pb"
	inventorySrv "github.com/Yoshikrit/inventory/internal/service/inventory"
)

func (c *InventoryController) GetProducts(ctx context.Context, _ *pb.GetProductsRequest) (*pb.ProductListResponse, error) {
	products, err := c.inventoryService.GetProducts(ctx)
	if err != nil {
		return nil, err
	}
	return parseProductListResponse(products), nil
}

func parseProductListResponse(resp inventorySrv.GetProductsResponse) *pb.ProductListResponse {
	var list []*pb.ProductResponse
	for _, p := range resp.Products {
		list = append(list, &pb.ProductResponse{
			ProductId:   p.ProductID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    uint64(p.Quantity),
		})
	}
	return &pb.ProductListResponse{Products: list}
}
