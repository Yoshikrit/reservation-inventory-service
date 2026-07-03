package mocks

import (
	"context"

	"github.com/Yoshikrit/inventory/internal/entity"
	"github.com/Yoshikrit/inventory/internal/pkg/apperror"

	"github.com/stretchr/testify/mock"
)

type ProductStockHistoryRepository struct {
	mock.Mock
}

func (m *ProductStockHistoryRepository) Create(ctx context.Context, history *entity.ProductStockHistory) *apperror.AppError {
	args := m.Called(ctx, history)
	if v := args.Get(0); v != nil {
		return v.(*apperror.AppError)
	}
	return nil
}
