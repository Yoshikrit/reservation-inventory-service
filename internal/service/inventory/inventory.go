package inventory

import (
	"context"
	"fmt"
	"time"

	"inventory/internal/pkg/apperror"
	historyRepo "inventory/internal/repository/producthistory"
	productRepo "inventory/internal/repository/product"

	"github.com/redis/go-redis/v9"
)

const (
	cacheKeyGetProducts = "service::GetProducts"
	cacheTTL            = 5 * time.Minute
)

func cacheKeyProductByID(productID string) string {
	return fmt.Sprintf("service::GetProductByID::%s", productID)
}

type InventoryService interface {
	CreateProduct(ctx context.Context, request *CreateProductRequest) *apperror.AppError
	GetProductByID(ctx context.Context, request *GetProductRequest) (GetProductResponse, *apperror.AppError)
	GetProducts(ctx context.Context) (GetProductsResponse, *apperror.AppError)
	CheckAndHold(ctx context.Context, request *CheckAndHoldRequest) (*CheckAndHoldResponse, *apperror.AppError)
	DeductStock(ctx context.Context, request *DeductStockRequest) *apperror.AppError
}

type TrManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type inventoryService struct {
	productRepo productRepo.ProductRepository
	historyRepo historyRepo.ProductStockHistoryRepository
	redis       *redis.Client
	trManager   TrManager
}

func NewInventoryService(
	productRepo productRepo.ProductRepository,
	historyRepo historyRepo.ProductStockHistoryRepository,
	rdb *redis.Client,
	trManager TrManager,
) InventoryService {
	return &inventoryService{
		productRepo: productRepo,
		historyRepo: historyRepo,
		redis:       rdb,
		trManager:   trManager,
	}
}
