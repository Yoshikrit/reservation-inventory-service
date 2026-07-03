package inventory

type CreateProductRequest struct {
	ProductID   string  `json:"product_id" validate:"required"`
	Name        string  `json:"name"       validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price"      validate:"required,gt=0"`
	Quantity    uint    `json:"quantity"   validate:"required,gt=0"`
}

type GetProductByIDRequest struct {
	ProductID string `uri:"product_id" validate:"required"`
}

type ProductResponse struct {
	ProductID   string  `json:"product_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    uint    `json:"quantity"`
}

type ProductListResponse struct {
	ProductList []ProductResponse `json:"products"`
}
