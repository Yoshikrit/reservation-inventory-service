package product

import (
	"context"

	"github.com/Yoshikrit/inventory/internal/entity"
	"github.com/Yoshikrit/inventory/internal/pkg/apperror"

	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(ctx context.Context, product *entity.Product) *apperror.AppError
	CreateBulk(ctx context.Context, products []entity.Product) *apperror.AppError
	Update(ctx context.Context, product *entity.Product) *apperror.AppError
	Patch(ctx context.Context, product *entity.Product) *apperror.AppError
	Find(ctx context.Context, filter *entity.Product) (*entity.Product, *apperror.AppError)
	Filter(ctx context.Context, filter *entity.Product, limit, offset int, isAsc bool) ([]entity.Product, *apperror.AppError)
	FindAll(ctx context.Context) ([]entity.Product, *apperror.AppError)
	Delete(ctx context.Context, filter *entity.Product) *apperror.AppError
	Count(ctx context.Context, filter *entity.Product) (int64, *apperror.AppError)

	// custom
	FindForUpdate(ctx context.Context, filter *entity.Product) (*entity.Product, *apperror.AppError)
	UpdateQuantity(ctx context.Context, productID string, quantity uint) *apperror.AppError
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db}
}
