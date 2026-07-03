package inventory_test

import (
	"context"
	"encoding/json"
	"testing"

	kafkaCtrl "inventory/internal/controller/kafka/v1/inventory"
	"inventory/internal/pkg/apperror"
	svc "inventory/internal/service/inventory"
	svcMocks "inventory/internal/service/inventory/mocks"

	kafka "github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newMsg(t *testing.T, payload any) kafka.Message {
	t.Helper()
	data, _ := json.Marshal(payload)
	return kafka.Message{Topic: "reservation-confirmed", Partition: 0, Offset: 1, Value: data}
}

// ── HandleConfirmedReservation ─────────────────────────────────────────────────

func TestKafka_HandleConfirmedReservation_Success(t *testing.T) {
	mockSvc := new(svcMocks.InventoryService)
	mockSvc.On("DeductStock", mock.Anything, mock.MatchedBy(func(r *svc.DeductStockRequest) bool {
		return r.ProductID == "prod-001" && r.Quantity == 3
	})).Return(nil)

	ctrl := kafkaCtrl.NewKafkaInventoryController(mockSvc)
	err := ctrl.HandleConfirmedReservation(context.Background(), newMsg(t, map[string]any{
		"product_id": "prod-001",
		"quantity":   3,
	}))

	assert.NoError(t, err)
	mockSvc.AssertExpectations(t)
}

func TestKafka_HandleConfirmedReservation_InvalidJSON(t *testing.T) {
	mockSvc := new(svcMocks.InventoryService)

	msg := kafka.Message{Value: []byte("not-json")}
	ctrl := kafkaCtrl.NewKafkaInventoryController(mockSvc)
	err := ctrl.HandleConfirmedReservation(context.Background(), msg)

	assert.Error(t, err)
	mockSvc.AssertNotCalled(t, "DeductStock")
}

func TestKafka_HandleConfirmedReservation_MissingFields(t *testing.T) {
	mockSvc := new(svcMocks.InventoryService)

	ctrl := kafkaCtrl.NewKafkaInventoryController(mockSvc)
	err := ctrl.HandleConfirmedReservation(context.Background(), newMsg(t, map[string]any{
		"product_id": "", // missing quantity, empty product_id
	}))

	assert.Error(t, err)
	mockSvc.AssertNotCalled(t, "DeductStock")
}

func TestKafka_HandleConfirmedReservation_ServiceError(t *testing.T) {
	mockSvc := new(svcMocks.InventoryService)
	mockSvc.On("DeductStock", mock.Anything, mock.Anything).
		Return(apperror.NewError(42200000, nil))

	ctrl := kafkaCtrl.NewKafkaInventoryController(mockSvc)
	err := ctrl.HandleConfirmedReservation(context.Background(), newMsg(t, map[string]any{
		"product_id": "prod-001",
		"quantity":   999,
	}))

	assert.Error(t, err)
}
