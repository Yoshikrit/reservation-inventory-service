package inventory

import (
	"context"
	"errors"

	"inventory/internal/entity"
	"inventory/internal/pkg/apperror"
	"inventory/internal/service/constant"

	"github.com/rs/zerolog/log"
)

func (s *inventoryService) DeductStock(ctx context.Context, req *DeductStockRequest) *apperror.AppError {
	if err := s.trManager.Do(ctx, func(ctx context.Context) error {
		product, appErr := s.productRepo.FindForUpdate(ctx, &entity.Product{ProductID: req.ProductID})
		if appErr != nil {
			return appErr
		}

		if product.Quantity < req.Quantity {
			return apperror.NewError(42200000, errors.New("insufficient stock"))
		}

		newQty := product.Quantity - req.Quantity
		if appErr := s.productRepo.UpdateQuantity(ctx, req.ProductID, newQty); appErr != nil {
			return appErr
		}
		if appErr := s.historyRepo.Create(ctx, &entity.ProductStockHistory{
			ProductID: req.ProductID,
			OldQty:    product.Quantity,
			NewQty:    newQty,
			Delta:     -int(req.Quantity),
			Reason:    constant.StockReasonReservationConfirmed,
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
	if err := s.redis.Del(ctx, cacheKeyProductByID(req.ProductID), cacheKeyGetProducts).Err(); err != nil {
		log.Warn().Err(err).Msg("redis del failed for product cache")
	}
	return nil
}
