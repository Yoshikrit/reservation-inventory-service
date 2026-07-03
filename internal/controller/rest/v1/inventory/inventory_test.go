package inventory_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"inventory/config"
	"inventory/internal/controller/rest"
	ctrlRest "inventory/internal/controller/rest/v1/inventory"
	"inventory/internal/pkg/apperror"
	svc "inventory/internal/service/inventory"
	svcMocks "inventory/internal/service/inventory/mocks"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newApp(svcMock *svcMocks.InventoryService) *fiber.App {
	app := fiber.New(config.NewRestConfig(rest.ErrorHandler()))
	ctrl := ctrlRest.NewInventoryController(svcMock)
	ctrl.RegisterRoutes(app.Group("/inventories"))
	return app
}

// ── CreateProduct ──────────────────────────────────────────────────────────────

func TestREST_CreateProduct_Success(t *testing.T) {
	mockSvc := new(svcMocks.InventoryService)
	mockSvc.On("CreateProduct", mock.Anything, mock.MatchedBy(func(r *svc.CreateProductRequest) bool {
		return r.ProductID == "prod-001" && r.Price == 999.99 && r.Quantity == 10
	})).Return(nil)

	body, _ := json.Marshal(map[string]any{
		"product_id": "prod-001",
		"name":       "Laptop",
		"price":      999.99,
		"quantity":   10,
	})
	req := httptest.NewRequest("POST", "/inventories", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := newApp(mockSvc).Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	mockSvc.AssertExpectations(t)
}

func TestREST_CreateProduct_ValidationError(t *testing.T) {
	mockSvc := new(svcMocks.InventoryService)

	body, _ := json.Marshal(map[string]any{"name": "Laptop"}) // missing product_id, price, quantity
	req := httptest.NewRequest("POST", "/inventories", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := newApp(mockSvc).Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	mockSvc.AssertNotCalled(t, "CreateProduct")
}

func TestREST_CreateProduct_ServiceError(t *testing.T) {
	mockSvc := new(svcMocks.InventoryService)
	mockSvc.On("CreateProduct", mock.Anything, mock.Anything).
		Return(apperror.NewError(40900000, nil, "product", "prod-001"))

	body, _ := json.Marshal(map[string]any{
		"product_id": "prod-001",
		"name":       "Laptop",
		"price":      999.99,
		"quantity":   10,
	})
	req := httptest.NewRequest("POST", "/inventories", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := newApp(mockSvc).Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusConflict, resp.StatusCode)
}

// ── GetProductByID ─────────────────────────────────────────────────────────────

func TestREST_GetProductByID_Success(t *testing.T) {
	mockSvc := new(svcMocks.InventoryService)
	mockSvc.On("GetProductByID", mock.Anything, mock.MatchedBy(func(r *svc.GetProductRequest) bool {
		return r.ProductID == "prod-001"
	})).Return(svc.GetProductResponse{
		ProductID: "prod-001",
		Name:      "Laptop",
		Price:     999.99,
		Quantity:  10,
	}, nil)

	req := httptest.NewRequest("GET", "/inventories/prod-001", nil)
	resp, err := newApp(mockSvc).Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var body map[string]any
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, "prod-001", body["product_id"])
	assert.Equal(t, "Laptop", body["name"])
}

func TestREST_GetProductByID_NotFound(t *testing.T) {
	mockSvc := new(svcMocks.InventoryService)
	mockSvc.On("GetProductByID", mock.Anything, mock.Anything).
		Return(svc.GetProductResponse{}, apperror.NewError(40400000, nil, "product", "prod-x"))

	req := httptest.NewRequest("GET", "/inventories/prod-x", nil)
	resp, err := newApp(mockSvc).Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

// ── GetProducts ────────────────────────────────────────────────────────────────

func TestREST_GetProducts_Success(t *testing.T) {
	mockSvc := new(svcMocks.InventoryService)
	mockSvc.On("GetProducts", mock.Anything).Return(svc.GetProductsResponse{
		Products: []svc.GetProductResponse{
			{ProductID: "prod-001", Name: "Laptop"},
			{ProductID: "prod-002", Name: "Mouse"},
		},
	}, nil)

	req := httptest.NewRequest("GET", "/inventories", nil)
	resp, err := newApp(mockSvc).Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var body map[string]any
	json.NewDecoder(resp.Body).Decode(&body)
	products := body["products"].([]any)
	assert.Len(t, products, 2)
}

func TestREST_GetProducts_InternalError(t *testing.T) {
	mockSvc := new(svcMocks.InventoryService)
	mockSvc.On("GetProducts", mock.Anything).
		Return(svc.GetProductsResponse{}, apperror.NewError(50000000, nil))

	req := httptest.NewRequest("GET", "/inventories", nil)
	resp, err := newApp(mockSvc).Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}
