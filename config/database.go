package config

import (
	"github.com/Yoshikrit/inventory/internal/entity"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatabaseConfig struct {
	DatabaseUrl string `env:"DATABASE_URL,required"`
}

func InitDatabase(databaseUrl string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseUrl), &gorm.Config{
		TranslateError: true,
		Logger:         logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func MigrateDatabase(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.Product{},
		&entity.ProductStockHistory{},
	)
}
