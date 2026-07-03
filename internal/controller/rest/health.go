package rest

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"gorm.io/gorm"
)

func newHealthRouter(router *fiber.App, db *gorm.DB) {
	router.Get(healthcheck.LivenessEndpoint, healthcheck.New())
	router.Get(healthcheck.ReadinessEndpoint, healthcheck.New(healthcheck.Config{
		Probe: func(c fiber.Ctx) bool {
			sqlDB, err := db.DB()
			if err != nil {
				return false
			}
			return sqlDB.PingContext(c.Context()) == nil
		},
	}))
}
