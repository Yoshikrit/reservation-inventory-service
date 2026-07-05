package metrics

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	}, []string{"method", "path", "status"})

	httpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP request duration in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path"})
)

func Metrics() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()
		err := c.Next()

		httpRequestsTotal.WithLabelValues(
			c.Method(),
			c.Route().Path,
			strconv.Itoa(c.Response().StatusCode()),
		).Inc()

		httpRequestDuration.WithLabelValues(
			c.Method(),
			c.Route().Path,
		).Observe(time.Since(start).Seconds())

		return err
	}
}
