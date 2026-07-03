package inventory

import (
	"context"

	"github.com/Yoshikrit/inventory/internal/pkg/apperror"
	"github.com/Yoshikrit/inventory/internal/pkg/validator"
	"github.com/Yoshikrit/inventory/internal/pkg/json"
	inventorySrv "github.com/Yoshikrit/inventory/internal/service/inventory"

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
