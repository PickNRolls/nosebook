package middleware

import (
	"nosebook/src/lib/clock"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var TotalRequestsCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of requests coming into HTTP server, excluding /metrics handler",
	},
	[]string{"path"},
)

var InProgressRequestsGauge = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "http_requests_in_progress",
		Help: "Number of \"in progress\" requests coming into HTTP server,  which are still handling, excluding /metrics handler",
	},
	[]string{"path"},
)

var ResponseElapsedHist = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Buckets: prometheus.ExponentialBuckets(0.001, 2, 12),
		Name:    "http_response_elapsed_seconds",
		Help:    "Total elapsed time of HTTP request",
	},
	[]string{"path"},
)

func NewRequestMetrics() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
    path := ctx.FullPath()

		TotalRequestsCounter.WithLabelValues(path).Inc()
		InProgressRequestsGauge.WithLabelValues(path).Inc()
		before := clock.Now()

		ctx.Next()

		after := clock.Now()
		InProgressRequestsGauge.WithLabelValues(path).Dec()
		ResponseElapsedHist.WithLabelValues(path).Observe(float64(after.Sub(before).Seconds()))
	}
}
