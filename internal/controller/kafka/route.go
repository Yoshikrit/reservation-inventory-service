package kafka

import (
	"inventory/config"
	"inventory/internal/controller/kafka/middleware"
	kafkaidempotency "inventory/internal/controller/kafka/middleware/idempotency"
	inventoryCtrl "inventory/internal/controller/kafka/v1/inventory"
	historyRepo "inventory/internal/repository/producthistory"
	productRepo "inventory/internal/repository/product"
	inventorySrv "inventory/internal/service/inventory"

	gormtrm "github.com/avito-tech/go-transaction-manager/drivers/gorm/v2"
	trm "github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type consumerConfig struct {
	topic   string
	groupID string
	handler middleware.Handler
}

func NewKafkaRouter(cfg config.KafkaConfig, db *gorm.DB, rdb *redis.Client) []consumerConfig {
	trManager, err := trm.New(gormtrm.NewDefaultFactory(db))
	if err != nil {
		log.Fatal().Err(err).Msg("kafka: failed to create transaction manager")
	}

	repo := productRepo.NewProductRepository(db)
	history := historyRepo.NewProductStockHistoryRepository(db)
	svc := inventorySrv.NewInventoryService(repo, history, rdb, trManager)
	ctrl := inventoryCtrl.NewKafkaInventoryController(svc)

	return []consumerConfig{
		{
			topic:   cfg.TopicConfirmedReservation,
			groupID: cfg.GroupInventoryConfirmedReservation,
			handler: middleware.Chain(
				ctrl.HandleConfirmedReservation,
				kafkaidempotency.Idempotency(rdb),
			),
		},
	}
}
