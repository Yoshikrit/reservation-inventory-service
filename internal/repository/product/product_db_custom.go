package product

import (
	"context"
	"errors"
	"github.com/Yoshikrit/inventory/internal/entity"
	"github.com/Yoshikrit/inventory/internal/pkg/apperror"

	gormtrm "github.com/avito-tech/go-transaction-manager/drivers/gorm/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *productRepository) FindForUpdate(ctx context.Context, filter *entity.Product) (*entity.Product, *apperror.AppError) {
	db := gormtrm.DefaultCtxGetter.DefaultTrOrDB(ctx, r.db)
	primaryFilter, appErr := r.buildPrimaryFilter(filter)
	if appErr != nil {
		return nil, appErr
	}

	var product entity.Product
	err := db.WithContext(ctx).Where(primaryFilter).Clauses(clause.Locking{Strength: "UPDATE"}).First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewError(40400000, err, "product", primaryFilter["product_id"])
		}
		return nil, apperror.NewError(50000000, err)
	}
	return &product, nil
}

func (r *productRepository) UpdateQuantity(ctx context.Context, productID string, quantity uint) *apperror.AppError {
	db := gormtrm.DefaultCtxGetter.DefaultTrOrDB(ctx, r.db)
	if err := db.WithContext(ctx).Model(&entity.Product{}).Where("product_id = ?", productID).Update("quantity", quantity).Error; err != nil {
		return apperror.NewError(50000000, err)
	}
	return nil
}
