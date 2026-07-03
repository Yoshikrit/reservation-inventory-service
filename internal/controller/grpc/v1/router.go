package v1

import (
	inventoryCtrl "github.com/Yoshikrit/inventory/internal/controller/grpc/v1/inventory"
	"github.com/Yoshikrit/inventory/internal/controller/grpc/v1/inventory/pb"
	historyRepo "github.com/Yoshikrit/inventory/internal/repository/producthistory"
	productRepo "github.com/Yoshikrit/inventory/internal/repository/product"
	inventorySrv "github.com/Yoshikrit/inventory/internal/service/inventory"

	gormtrm "github.com/avito-tech/go-transaction-manager/drivers/gorm/v2"
	trm "github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

func NewGRPCRouter(server *grpc.Server, db *gorm.DB, rdb *redis.Client) {
	trManager, err := trm.New(gormtrm.NewDefaultFactory(db))
	if err != nil {
		log.Fatal().Err(err).Msg("grpc: failed to create transaction manager")
	}

	productRepo := productRepo.NewProductRepository(db)
	historyRepo := historyRepo.NewProductStockHistoryRepository(db)
	inventorySvc := inventorySrv.NewInventoryService(productRepo, historyRepo, rdb, trManager)
	inventoryCtrl := inventoryCtrl.NewInventoryController(inventorySvc)

	pb.RegisterInventoryServiceServer(server, inventoryCtrl)
}
