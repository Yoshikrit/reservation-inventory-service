package inventory

import (
	"context"
	"errors"

	"github.com/Yoshikrit/inventory/internal/entity"
	"github.com/Yoshikrit/inventory/internal/pkg/apperror"
	"github.com/Yoshikrit/inventory/internal/service/constant"

	"github.com/rs/zerolog/log"
)

func (s *inventoryService) CreateProduct(ctx context.Context, product *CreateProductRequest) *apperror.AppError {
	if err := s.trManager.Do(ctx, func(ctx context.Context) error {
		productEntity := entity.Product{
			ProductID:   product.ProductID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Quantity:    product.Quantity,
		}
		if appErr := s.productRepo.Create(ctx, &productEntity); appErr != nil {
			return appErr
		}
		if appErr := s.historyRepo.Create(ctx, &entity.ProductStockHistory{
			ProductID: product.ProductID,
			OldQty:    0,
			NewQty:    product.Quantity,
			Delta:     int(product.Quantity),
			Reason:    constant.StockReasonProductCreated,
		}); appErr != nil {
			return appErr
		}
		return nil
	}); err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			return appErr
		}
		return apperror.NewError(50000000, err)
	}

	if err := s.redis.Del(ctx, cacheKeyGetProducts).Err(); err != nil {
		log.Warn().Err(err).Msg("redis del failed for products cache")
	}
	return nil
}
