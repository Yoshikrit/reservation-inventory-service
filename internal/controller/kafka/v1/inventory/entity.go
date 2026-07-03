package inventory

type ConfirmedReservationEvent struct {
	ProductID string `json:"product_id" validate:"required"`
	Quantity  uint   `json:"quantity"   validate:"required,gt=0"`
}
