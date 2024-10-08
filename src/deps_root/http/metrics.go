package roothttp

import (
	"nosebook/src/deps_root/http/middleware"
	"nosebook/src/deps_root/worker"

	"github.com/prometheus/client_golang/prometheus"
)

func registerMetrics() {
	prometheus.MustRegister(middleware.TotalRequestsCounter)
	prometheus.MustRegister(middleware.InProgressRequestsGauge)
	prometheus.MustRegister(middleware.ResponseElapsedHist)
	prometheus.MustRegister(InProgressWsConnectionsGauge)

	prometheus.MustRegister(middleware.DBIdleConnectionsGauge)
	prometheus.MustRegister(middleware.DBActiveConnectionsGauge)
	prometheus.MustRegister(middleware.DBOpenConnectionsGauge)
	prometheus.MustRegister(middleware.DBWaitCountGauge)
	prometheus.MustRegister(middleware.DBWaitSecondsGauge)

	prometheus.MustRegister(worker.FlushedBufferSize)
	prometheus.MustRegister(worker.FlushElapsedSeconds)
	prometheus.MustRegister(worker.SendElapsedSeconds)
}
