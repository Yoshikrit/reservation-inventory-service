package mocks

import (
	"context"

	"github.com/Yoshikrit/inventory/internal/pkg/apperror"
	svc "github.com/Yoshikrit/inventory/internal/service/inventory"

	"github.com/stretchr/testify/mock"
)

type InventoryService struct {
	mock.Mock
}

func (m *InventoryService) CreateProduct(ctx context.Context, request *svc.CreateProductRequest) *apperror.AppError {
	args := m.Called(ctx, request)
	if v := args.Get(0); v != nil {
		return v.(*apperror.AppError)
	}
	return nil
}

func (m *InventoryService) GetProductByID(ctx context.Context, request *svc.GetProductRequest) (svc.GetProductResponse, *apperror.AppError) {
	args := m.Called(ctx, request)
	if v := args.Get(1); v != nil {
		return svc.GetProductResponse{}, v.(*apperror.AppError)
	}
	return args.Get(0).(svc.GetProductResponse), nil
}

func (m *InventoryService) GetProducts(ctx context.Context) (svc.GetProductsResponse, *apperror.AppError) {
	args := m.Called(ctx)
	if v := args.Get(1); v != nil {
		return svc.GetProductsResponse{}, v.(*apperror.AppError)
	}
	return args.Get(0).(svc.GetProductsResponse), nil
}

func (m *InventoryService) CheckAndHold(ctx context.Context, request *svc.CheckAndHoldRequest) (*svc.CheckAndHoldResponse, *apperror.AppError) {
	args := m.Called(ctx, request)
	if v := args.Get(1); v != nil {
		return nil, v.(*apperror.AppError)
	}
	return args.Get(0).(*svc.CheckAndHoldResponse), nil
}

func (m *InventoryService) DeductStock(ctx context.Context, request *svc.DeductStockRequest) *apperror.AppError {
	args := m.Called(ctx, request)
	if v := args.Get(0); v != nil {
		return v.(*apperror.AppError)
	}
	return nil
}
