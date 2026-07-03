package inventory

import (
	"context"

	"inventory/internal/pkg/apperror"
	"inventory/internal/pkg/validator"
	"inventory/internal/pkg/json"
	inventorySrv "inventory/internal/service/inventory"

	kafka "github.com/segmentio/kafka-go"
)

var v = validator.New()

func (c *KafkaInventoryController) HandleConfirmedReservation(ctx context.Context, msg kafka.Message) error {
	var request ConfirmedReservationEvent
	if err := json.Unmarshal(msg.Value, &request); err != nil {
		return apperror.NewError(40000000, err)
	}

	if err := v.Validate(request); err != nil {
		return apperror.NewError(40000001, err)
	}

	if appErr := c.inventoryService.DeductStock(ctx, &inventorySrv.DeductStockRequest{
		ProductID: request.ProductID,
		Quantity:  request.Quantity,
	}); appErr != nil {
		return appErr
	}
	return nil
}
