package inventory_test

import (
	"context"
	"testing"

	"inventory/internal/entity"
	"inventory/internal/pkg/apperror"
	productMocks "inventory/internal/repository/product/mocks"
	historyMocks "inventory/internal/repository/producthistory/mocks"
	"inventory/internal/service/constant"
	svc "inventory/internal/service/inventory"
	svcMocks "inventory/internal/service/inventory/mocks"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newRedis(t *testing.T) *redis.Client {
	t.Helper()
	mr := miniredis.RunT(t)
	return redis.NewClient(&redis.Options{Addr: mr.Addr()})
}

func trPassthrough(trMgr *svcMocks.TrManager) {
	trMgr.On("Do", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(context.Context) error)
			fn(context.Background())
		}).
		Return(nil)
}

// ── CreateProduct ──────────────────────────────────────────────────────────────

func TestCreateProduct_Success(t *testing.T) {
	productRepo := new(productMocks.ProductRepository)
	historyRepo := new(historyMocks.ProductStockHistoryRepository)
	trMgr := new(svcMocks.TrManager)
	rdb := newRedis(t)

	trPassthrough(trMgr)
	productRepo.On("Create", mock.Anything, mock.MatchedBy(func(p *entity.Product) bool {
		return p.ProductID == "prod-001" && p.Quantity == 10
	})).Return(nil)
	historyRepo.On("Create", mock.Anything, mock.MatchedBy(func(h *entity.ProductStockHistory) bool {
		return h.ProductID == "prod-001" && h.Delta == 10 && h.Reason == constant.StockReasonProductCreated
	})).Return(nil)

	service := svc.NewInventoryService(productRepo, historyRepo, rdb, trMgr)
	appErr := service.CreateProduct(context.Background(), &svc.CreateProductRequest{
		ProductID: "prod-001",
		Name:      "Laptop",
		Price:     999.99,
		Quantity:  10,
	})

	assert.Nil(t, appErr)
	productRepo.AssertExpectations(t)
	historyRepo.AssertExpectations(t)
	trMgr.AssertExpectations(t)
}

func TestCreateProduct_ProductAlreadyExists(t *testing.T) {
	productRepo := new(productMocks.ProductRepository)
	historyRepo := new(historyMocks.ProductStockHistoryRepository)
	trMgr := new(svcMocks.TrManager)
	rdb := newRedis(t)

	conflictErr := apperror.NewError(40900000, nil, "product", "prod-001")
	trMgr.On("Do", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(context.Context) error)
			fn(context.Background())
		}).
		Return(conflictErr)
	productRepo.On("Create", mock.Anything, mock.Anything).Return(conflictErr)

	service := svc.NewInventoryService(productRepo, historyRepo, rdb, trMgr)
	appErr := service.CreateProduct(context.Background(), &svc.CreateProductRequest{
		ProductID: "prod-001",
		Name:      "Laptop",
		Price:     999.99,
		Quantity:  10,
	})

	assert.NotNil(t, appErr)
	assert.Equal(t, apperror.CategoryConflict, appErr.Category)
	historyRepo.AssertNotCalled(t, "Create")
}

func TestCreateProduct_HistoryFails_RollsBack(t *testing.T) {
	productRepo := new(productMocks.ProductRepository)
	historyRepo := new(historyMocks.ProductStockHistoryRepository)
	trMgr := new(svcMocks.TrManager)
	rdb := newRedis(t)

	internalErr := apperror.NewError(50000000, nil)
	trMgr.On("Do", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(context.Context) error)
			fn(context.Background())
		}).
		Return(internalErr)
	productRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
	historyRepo.On("Create", mock.Anything, mock.Anything).Return(internalErr)

	service := svc.NewInventoryService(productRepo, historyRepo, rdb, trMgr)
	appErr := service.CreateProduct(context.Background(), &svc.CreateProductRequest{
		ProductID: "prod-001",
		Name:      "Laptop",
		Price:     999.99,
		Quantity:  10,
	})

	assert.NotNil(t, appErr)
	assert.Equal(t, apperror.CategoryInternal, appErr.Category)
}

// ── GetProductByID ─────────────────────────────────────────────────────────────

func TestGetProductByID_CacheHit(t *testing.T) {
	productRepo := new(productMocks.ProductRepository)
	historyRepo := new(historyMocks.ProductStockHistoryRepository)
	trMgr := new(svcMocks.TrManager)
	rdb := newRedis(t)

	ctx := context.Background()
	rdb.Set(ctx, "service::GetProductByID::prod-001",
		`{"ProductID":"prod-001","Name":"Laptop","Description":"","Price":999.99,"Quantity":10}`,
		0)

	service := svc.NewInventoryService(productRepo, historyRepo, rdb, trMgr)
	resp, appErr := service.GetProductByID(ctx, &svc.GetProductRequest{ProductID: "prod-001"})

	assert.Nil(t, appErr)
	assert.Equal(t, "prod-001", resp.ProductID)
	assert.Equal(t, "Laptop", resp.Name)
	productRepo.AssertNotCalled(t, "Find")
}

func TestGetProductByID_CacheMiss_FallsToDb(t *testing.T) {
	productRepo := new(productMocks.ProductRepository)
	historyRepo := new(historyMocks.ProductStockHistoryRepository)
	trMgr := new(svcMocks.TrManager)
	rdb := newRedis(t)

	productRepo.On("Find", mock.Anything, mock.MatchedBy(func(p *entity.Product) bool {
		return p.ProductID == "prod-001"
	})).Return(&entity.Product{
		ProductID: "prod-001",
		Name:      "Laptop",
		Price:     999.99,
		Quantity:  10,
	}, nil)

	service := svc.NewInventoryService(productRepo, historyRepo, rdb, trMgr)
	resp, appErr := service.GetProductByID(context.Background(), &svc.GetProductRequest{ProductID: "prod-001"})

	assert.Nil(t, appErr)
	assert.Equal(t, "prod-001", resp.ProductID)
	productRepo.AssertExpectations(t)
}

func TestGetProductByID_NotFound(t *testing.T) {
	productRepo := new(productMocks.ProductRepository)
	historyRepo := new(historyMocks.ProductStockHistoryRepository)
	trMgr := new(svcMocks.TrManager)
	rdb := newRedis(t)

	notFoundErr := apperror.NewError(40400000, nil, "product", "no-such-product")
	productRepo.On("Find", mock.Anything, mock.Anything).Return((*entity.Product)(nil), notFoundErr)

	service := svc.NewInventoryService(productRepo, historyRepo, rdb, trMgr)
	_, appErr := service.GetProductByID(context.Background(), &svc.GetProductRequest{ProductID: "no-such-product"})

	assert.NotNil(t, appErr)
	assert.Equal(t, apperror.CategoryNotFound, appErr.Category)
}

// ── GetProducts ────────────────────────────────────────────────────────────────

func TestGetProducts_CacheMiss_ReturnsFromDb(t *testing.T) {
	productRepo := new(productMocks.ProductRepository)
	historyRepo := new(historyMocks.ProductStockHistoryRepository)
	trMgr := new(svcMocks.TrManager)
	rdb := newRedis(t)

	productRepo.On("FindAll", mock.Anything).Return([]entity.Product{
		{ProductID: "prod-001", Name: "Laptop", Price: 999.99, Quantity: 10},
		{ProductID: "prod-002", Name: "Mouse", Price: 29.99, Quantity: 50},
	}, nil)

	service := svc.NewInventoryService(productRepo, historyRepo, rdb, trMgr)
	resp, appErr := service.GetProducts(context.Background())

	assert.Nil(t, appErr)
	assert.Len(t, resp.Products, 2)
	assert.Equal(t, "prod-001", resp.Products[0].ProductID)
	productRepo.AssertExpectations(t)
}

func TestGetProducts_EmptyList(t *testing.T) {
	productRepo := new(productMocks.ProductRepository)
	historyRepo := new(historyMocks.ProductStockHistoryRepository)
	trMgr := new(svcMocks.TrManager)
	rdb := newRedis(t)

	productRepo.On("FindAll", mock.Anything).Return([]entity.Product{}, nil)

	service := svc.NewInventoryService(productRepo, historyRepo, rdb, trMgr)
	resp, appErr := service.GetProducts(context.Background())

	assert.Nil(t, appErr)
	assert.Empty(t, resp.Products)
}

// ── CheckAndHold ───────────────────────────────────────────────────────────────

func TestCheckAndHold_Sufficient(t *testing.T) {
	productRepo := new(productMocks.ProductRepository)
	historyRepo := new(historyMocks.ProductStockHistoryRepository)
	trMgr := new(svcMocks.TrManager)
	rdb := newRedis(t)

	trPassthrough(trMgr)
	productRepo.On("FindForUpdate", mock.Anything, mock.MatchedBy(func(p *entity.Product) bool {
		return p.ProductID == "prod-001"
	})).Return(&entity.Product{
		ProductID: "prod-001",
		Name:      "Laptop",
		Price:     999.99,
		Quantity:  10,
	}, nil)

	service := svc.NewInventoryService(productRepo, historyRepo, rdb, trMgr)
	resp, appErr := service.CheckAndHold(context.Background(), &svc.CheckAndHoldRequest{
		ProductID: "prod-001",
		Quantity:  5,
	})

	assert.Nil(t, appErr)
	assert.True(t, resp.Available)
	assert.Equal(t, float64(999.99), resp.Price)
}

func TestCheckAndHold_Insufficient(t *testing.T) {
	productRepo := new(productMocks.ProductRepository)
	historyRepo := new(historyMocks.ProductStockHistoryRepository)
	trMgr := new(svcMocks.TrManager)
	rdb := newRedis(t)

	trPassthrough(trMgr)
	productRepo.On("FindForUpdate", mock.Anything, mock.Anything).Return(&entity.Product{
		ProductID: "prod-001",
		Quantity:  3,
	}, nil)

	service := svc.NewInventoryService(productRepo, historyRepo, rdb, trMgr)
	resp, appErr := service.CheckAndHold(context.Background(), &svc.CheckAndHoldRequest{
		ProductID: "prod-001",
		Quantity:  10,
	})

	assert.Nil(t, appErr)
	assert.False(t, resp.Available)
}

// ── DeductStock ────────────────────────────────────────────────────────────────

func TestDeductStock_Success(t *testing.T) {
	productRepo := new(productMocks.ProductRepository)
	historyRepo := new(historyMocks.ProductStockHistoryRepository)
	trMgr := new(svcMocks.TrManager)
	rdb := newRedis(t)

	trPassthrough(trMgr)
	productRepo.On("FindForUpdate", mock.Anything, mock.MatchedBy(func(p *entity.Product) bool {
		return p.ProductID == "prod-001"
	})).Return(&entity.Product{ProductID: "prod-001", Quantity: 10}, nil)
	productRepo.On("UpdateQuantity", mock.Anything, "prod-001", uint(8)).Return(nil)
	historyRepo.On("Create", mock.Anything, mock.MatchedBy(func(h *entity.ProductStockHistory) bool {
		return h.ProductID == "prod-001" && h.Delta == -2 && h.Reason == constant.StockReasonReservationConfirmed
	})).Return(nil)

	service := svc.NewInventoryService(productRepo, historyRepo, rdb, trMgr)
	appErr := service.DeductStock(context.Background(), &svc.DeductStockRequest{
		ProductID: "prod-001",
		Quantity:  2,
	})

	assert.Nil(t, appErr)
	productRepo.AssertExpectations(t)
	historyRepo.AssertExpectations(t)
}

func TestDeductStock_InsufficientStock(t *testing.T) {
	productRepo := new(productMocks.ProductRepository)
	historyRepo := new(historyMocks.ProductStockHistoryRepository)
	trMgr := new(svcMocks.TrManager)
	rdb := newRedis(t)

	insufficientErr := apperror.NewError(42200000, nil)
	trMgr.On("Do", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(context.Context) error)
			fn(context.Background())
		}).
		Return(insufficientErr)
	productRepo.On("FindForUpdate", mock.Anything, mock.Anything).
		Return(&entity.Product{ProductID: "prod-001", Quantity: 1}, nil)

	service := svc.NewInventoryService(productRepo, historyRepo, rdb, trMgr)
	appErr := service.DeductStock(context.Background(), &svc.DeductStockRequest{
		ProductID: "prod-001",
		Quantity:  5,
	})

	assert.NotNil(t, appErr)
	assert.Equal(t, apperror.CategoryUnprocessable, appErr.Category)
	historyRepo.AssertNotCalled(t, "Create")
}
