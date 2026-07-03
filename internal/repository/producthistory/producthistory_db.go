package producthistory

import (
	"context"

	"github.com/Yoshikrit/inventory/internal/entity"
	"github.com/Yoshikrit/inventory/internal/pkg/apperror"

	gormtrm "github.com/avito-tech/go-transaction-manager/drivers/gorm/v2"
)

func (r *productStockHistoryRepository) Create(ctx context.Context, history *entity.ProductStockHistory) *apperror.AppError {
	db := gormtrm.DefaultCtxGetter.DefaultTrOrDB(ctx, r.db)
	if err := db.WithContext(ctx).Create(history).Error; err != nil {
		return apperror.NewError(50000000, err)
	}
	return nil
}
