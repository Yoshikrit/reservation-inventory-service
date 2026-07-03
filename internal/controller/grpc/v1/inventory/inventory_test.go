package inventory_test

import (
	"context"
	"testing"

	grpcCtrl "github.com/Yoshikrit/inventory/internal/controller/grpc/v1/inventory"
	"github.com/Yoshikrit/inventory/internal/controller/grpc/v1/inventory/pb"
	"github.com/Yoshikrit/inventory/internal/pkg/apperror"
	svc "github.com/Yoshikrit/inventory/internal/service/inventory"
	svcMocks "github.com/Yoshikrit/inventory/internal/service/inventory/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ── CreateProduct ──────────────────────────────────────────────────────────────

func TestGRPC_CreateProduct_Success(t *testing.T) {
	mockSvc := new(svcMocks.InventoryService)
	mockSvc.On("CreateProduct", mock.Anything, mock.MatchedBy(func(r *svc.CreateProductRequest) bool {
		return r.ProductID == "prod-001" && r.Name == "Laptop" && r.Quantity == 5
	})).Return(nil)

	ctrl := grpcCtrl.NewInventoryController(mockSvc)
	resp, err := ctrl.CreateProduct(context.Background(), &pb.CreateProductRequest{
		ProductId: "prod-001",
		Name:      "Laptop",
		Price:     999.99,
		Quantity:  5,
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	mockSvc.AssertExpectations(t)
}

func TestGRPC_CreateProduct_ServiceError(t *testing.T) {
	mockSvc := new(svcMocks.InventoryService)
	mockSvc.On("CreateProduct", mock.Anything, mock.Anything).
		Return(apperror.NewError(40900000, nil, "product", "prod-001"))

	ctrl := grpcCtrl.NewInventoryController(mockSvc)
	resp, err := ctrl.CreateProduct(context.Background(), &pb.CreateProductRequest{
		ProductId: "prod-001",
		Name:      "Laptop",
		Price:     999.99,
		Quantity:  5,
	})

	assert.Error(t, err)
	assert.Nil(t, resp)
}

// ── GetProductByID ─────────────────────────────────────────────────────────────

func TestGRPC_GetProductByID_Success(t *testing.T) {
	mockSvc := new(svcMocks.InventoryService)
	mockSvc.On("GetProductByID", mock.Anything, mock.MatchedBy(func(r *svc.GetProductRequest) bool {
		return r.ProductID == "prod-001"
	})).Return(svc.GetProductResponse{
		ProductID: "prod-001",
		Name:      "Laptop",
		Price:     999.99,
		Quantity:  10,
	}, nil)

	ctrl := grpcCtrl.NewInventoryController(mockSvc)
	resp, err := ctrl.GetProductByID(context.Background(), &pb.GetProductByIDRequest{ProductId: "prod-001"})

	assert.NoError(t, err)
	assert.Equal(t, "prod-001", resp.ProductId)
	assert.Equal(t, "Laptop", resp.Name)
	assert.Equal(t, float64(999.99), resp.Price)
}

func TestGRPC_GetProductByID_NotFound(t *testing.T) {
	mockSvc := new(svcMocks.InventoryService)
	mockSvc.On("GetProductByID", mock.Anything, mock.Anything).
		Return(svc.GetProductResponse{}, apperror.NewError(40400000, nil, "product", "no-such"))

	ctrl := grpcCtrl.NewInventoryController(mockSvc)
	resp, err := ctrl.GetProductByID(context.Background(), &pb.GetProductByIDRequest{ProductId: "no-such"})

	assert.Error(t, err)
	assert.Nil(t, resp)
}

// ── GetProducts ────────────────────────────────────────────────────────────────

func TestGRPC_GetProducts_Success(t *testing.T) {
	mockSvc := new(svcMocks.InventoryService)
	mockSvc.On("GetProducts", mock.Anything).Return(svc.GetProductsResponse{
		Products: []svc.GetProductResponse{
			{ProductID: "prod-001", Name: "Laptop", Price: 999.99, Quantity: 10},
		},
	}, nil)

	ctrl := grpcCtrl.NewInventoryController(mockSvc)
	resp, err := ctrl.GetProducts(context.Background(), &pb.GetProductsRequest{})

	assert.NoError(t, err)
	assert.Len(t, resp.Products, 1)
	assert.Equal(t, "prod-001", resp.Products[0].ProductId)
}

// ── CheckAndHold ───────────────────────────────────────────────────────────────

func TestGRPC_CheckAndHold_Available(t *testing.T) {
	mockSvc := new(svcMocks.InventoryService)
	mockSvc.On("CheckAndHold", mock.Anything, mock.MatchedBy(func(r *svc.CheckAndHoldRequest) bool {
		return r.ProductID == "prod-001" && r.Quantity == 5
	})).Return(&svc.CheckAndHoldResponse{
		Available: true,
		ProductID: "prod-001",
		Name:      "Laptop",
		Price:     999.99,
		Quantity:  10,
	}, nil)

	ctrl := grpcCtrl.NewInventoryController(mockSvc)
	resp, err := ctrl.CheckAndHold(context.Background(), &pb.CheckAndHoldRequest{
		ProductId: "prod-001",
		Quantity:  5,
	})

	assert.NoError(t, err)
	assert.True(t, resp.Available)
	assert.Equal(t, "prod-001", resp.Product.ProductId)
}

func TestGRPC_CheckAndHold_NotFound(t *testing.T) {
	mockSvc := new(svcMocks.InventoryService)
	mockSvc.On("CheckAndHold", mock.Anything, mock.Anything).
		Return((*svc.CheckAndHoldResponse)(nil), apperror.NewError(40400000, nil, "product", "no-such"))

	ctrl := grpcCtrl.NewInventoryController(mockSvc)
	resp, err := ctrl.CheckAndHold(context.Background(), &pb.CheckAndHoldRequest{
		ProductId: "no-such",
		Quantity:  1,
	})

	assert.Error(t, err)
	assert.Nil(t, resp)
}
