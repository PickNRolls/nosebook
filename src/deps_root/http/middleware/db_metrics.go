package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
)

var DBIdleConnectionsGauge = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "db_idle_connections",
		Help: "Current number of idle connections",
	},
)

var DBActiveConnectionsGauge = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "db_active_connections",
		Help: "Current number of active connections",
	},
)

var DBWaitSecondsGauge = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "db_wait_total_seconds",
		Help: "Total time blocked waiting for new connection",
	},
)

var DBWaitCountGauge = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "db_wait_count",
		Help: "Total count of blocked waiters for new connection",
	},
)

var DBOpenConnectionsGauge = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "db_open_connections",
		Help: "Current number of open connections",
	},
)

func NewDbMetrics(db *sqlx.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		defer ctx.Next()

		stats := db.Stats()
		DBIdleConnectionsGauge.Set(float64(stats.Idle))
		DBActiveConnectionsGauge.Set(float64(stats.InUse))
		DBWaitSecondsGauge.Set(stats.WaitDuration.Seconds())
		DBWaitCountGauge.Set(float64(stats.WaitCount))
		DBOpenConnectionsGauge.Set(float64(stats.OpenConnections))
	}
}
