package rest

import (
	"inventory/internal/controller/rest/middleware"
	inventoryCtrl "inventory/internal/controller/rest/v1/inventory"
	historyRepo "inventory/internal/repository/producthistory"
	productRepo "inventory/internal/repository/product"
	inventorySrv "inventory/internal/service/inventory"

	gormtrm "github.com/avito-tech/go-transaction-manager/drivers/gorm/v2"
	trm "github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func NewRestRouter(router *fiber.App, db *gorm.DB, rdb *redis.Client) {
	for _, m := range middleware.NewMiddleware() {
		router.Use(m)
	}

	newHealthRouter(router, db)

	v1 := router.Group("/api/v1")
	newRoute(v1, db, rdb)
}

func newRoute(v1 fiber.Router, db *gorm.DB, rdb *redis.Client) {
	trManager, err := trm.New(gormtrm.NewDefaultFactory(db))
	if err != nil {
		log.Fatal().Err(err).Msg("rest: failed to create transaction manager")
	}

	productRepo := productRepo.NewProductRepository(db)
	historyRepo := historyRepo.NewProductStockHistoryRepository(db)
	inventorySvc := inventorySrv.NewInventoryService(productRepo, historyRepo, rdb, trManager)
	inventoryCtrl := inventoryCtrl.NewInventoryController(inventorySvc)

	inventoryRoute := v1.Group("/inventories")
	inventoryCtrl.RegisterRoutes(inventoryRoute)
}
