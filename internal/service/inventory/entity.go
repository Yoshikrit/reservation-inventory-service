package inventory

type CreateProductRequest struct {
	ProductID   string
	Name        string
	Description string
	Price       float64
	Quantity    uint
}

type GetProductRequest struct {
	ProductID string
}

type GetProductResponse struct {
	ProductID   string
	Name        string
	Description string
	Price       float64
	Quantity    uint
}

type GetProductsResponse struct {
	Products []GetProductResponse
}

type CheckAndHoldRequest struct {
	ProductID string
	Quantity  uint
}

type CheckAndHoldResponse struct {
	Available   bool
	ProductID   string
	Name        string
	Description string
	Price       float64
	Quantity    uint
}

type DeductStockRequest struct {
	ProductID string
	Quantity  uint
}
