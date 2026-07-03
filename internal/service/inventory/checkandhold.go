package inventory

import (
	"context"
	"errors"

	"github.com/Yoshikrit/inventory/internal/entity"
	"github.com/Yoshikrit/inventory/internal/pkg/apperror"
)

func (s *inventoryService) CheckAndHold(ctx context.Context, req *CheckAndHoldRequest) (*CheckAndHoldResponse, *apperror.AppError) {
	var result *CheckAndHoldResponse

	if err := s.trManager.Do(ctx, func(ctx context.Context) error {
		product, appErr := s.productRepo.FindForUpdate(ctx, &entity.Product{ProductID: req.ProductID})
		if appErr != nil {
			return appErr
		}

		result = &CheckAndHoldResponse{
			Available:   product.Quantity >= req.Quantity,
			ProductID:   product.ProductID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Quantity:    product.Quantity,
		}
		return nil
	}); err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			return nil, appErr
		}
		return nil, apperror.NewError(50000000, err)
	}

	return result, nil
}
