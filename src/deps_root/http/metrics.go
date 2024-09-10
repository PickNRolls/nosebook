package roothttp

import (
	"nosebook/src/deps_root/http/middleware"

	"github.com/prometheus/client_golang/prometheus"
)

func registerMetrics() {
	prometheus.MustRegister(PingCounter)
	prometheus.MustRegister(middleware.TotalRequestsCounter)
	prometheus.MustRegister(middleware.InProgressRequestsGauge)
	prometheus.MustRegister(middleware.ResponseElapsedHist)
}
