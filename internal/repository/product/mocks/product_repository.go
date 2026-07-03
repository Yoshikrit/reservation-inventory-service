package mocks

import (
	"context"

	"inventory/internal/entity"
	"inventory/internal/pkg/apperror"

	"github.com/stretchr/testify/mock"
)

type ProductRepository struct {
	mock.Mock
}

func (m *ProductRepository) Create(ctx context.Context, product *entity.Product) *apperror.AppError {
	args := m.Called(ctx, product)
	if v := args.Get(0); v != nil {
		return v.(*apperror.AppError)
	}
	return nil
}

func (m *ProductRepository) CreateBulk(ctx context.Context, products []entity.Product) *apperror.AppError {
	args := m.Called(ctx, products)
	if v := args.Get(0); v != nil {
		return v.(*apperror.AppError)
	}
	return nil
}

func (m *ProductRepository) Update(ctx context.Context, product *entity.Product) *apperror.AppError {
	args := m.Called(ctx, product)
	if v := args.Get(0); v != nil {
		return v.(*apperror.AppError)
	}
	return nil
}

func (m *ProductRepository) Patch(ctx context.Context, product *entity.Product) *apperror.AppError {
	args := m.Called(ctx, product)
	if v := args.Get(0); v != nil {
		return v.(*apperror.AppError)
	}
	return nil
}

func (m *ProductRepository) Find(ctx context.Context, filter *entity.Product) (*entity.Product, *apperror.AppError) {
	args := m.Called(ctx, filter)
	if v := args.Get(1); v != nil {
		return nil, v.(*apperror.AppError)
	}
	return args.Get(0).(*entity.Product), nil
}

func (m *ProductRepository) Filter(ctx context.Context, filter *entity.Product, limit, offset int, isAsc bool) ([]entity.Product, *apperror.AppError) {
	args := m.Called(ctx, filter, limit, offset, isAsc)
	if v := args.Get(1); v != nil {
		return nil, v.(*apperror.AppError)
	}
	return args.Get(0).([]entity.Product), nil
}

func (m *ProductRepository) FindAll(ctx context.Context) ([]entity.Product, *apperror.AppError) {
	args := m.Called(ctx)
	if v := args.Get(1); v != nil {
		return nil, v.(*apperror.AppError)
	}
	return args.Get(0).([]entity.Product), nil
}

func (m *ProductRepository) Delete(ctx context.Context, filter *entity.Product) *apperror.AppError {
	args := m.Called(ctx, filter)
	if v := args.Get(0); v != nil {
		return v.(*apperror.AppError)
	}
	return nil
}

func (m *ProductRepository) Count(ctx context.Context, filter *entity.Product) (int64, *apperror.AppError) {
	args := m.Called(ctx, filter)
	if v := args.Get(1); v != nil {
		return 0, v.(*apperror.AppError)
	}
	return args.Get(0).(int64), nil
}

func (m *ProductRepository) FindForUpdate(ctx context.Context, filter *entity.Product) (*entity.Product, *apperror.AppError) {
	args := m.Called(ctx, filter)
	if v := args.Get(1); v != nil {
		return nil, v.(*apperror.AppError)
	}
	return args.Get(0).(*entity.Product), nil
}

func (m *ProductRepository) UpdateQuantity(ctx context.Context, productID string, quantity uint) *apperror.AppError {
	args := m.Called(ctx, productID, quantity)
	if v := args.Get(0); v != nil {
		return v.(*apperror.AppError)
	}
	return nil
}
