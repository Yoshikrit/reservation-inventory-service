package producthistory

import (
	"context"

	"github.com/Yoshikrit/inventory/internal/entity"
	"github.com/Yoshikrit/inventory/internal/pkg/apperror"

	"gorm.io/gorm"
)

type ProductStockHistoryRepository interface {
	Create(ctx context.Context, history *entity.ProductStockHistory) *apperror.AppError
}

type productStockHistoryRepository struct {
	db *gorm.DB
}

func NewProductStockHistoryRepository(db *gorm.DB) ProductStockHistoryRepository {
	return &productStockHistoryRepository{db}
}
