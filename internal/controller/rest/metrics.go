package rest

import (
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/gofiber/fiber/v3"
)

func newMetricsRouter(router *fiber.App) {
	router.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
}
