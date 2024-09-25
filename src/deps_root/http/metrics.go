package roothttp

import (
	"nosebook/src/deps_root/http/middleware"
	repos "nosebook/src/infra/postgres/repositories"

	"github.com/prometheus/client_golang/prometheus"
)

func registerMetrics() {
	prometheus.MustRegister(middleware.TotalRequestsCounter)
	prometheus.MustRegister(middleware.InProgressRequestsGauge)
	prometheus.MustRegister(middleware.ResponseElapsedHist)
	prometheus.MustRegister(InProgressWsConnectionsGauge)

	prometheus.MustRegister(repos.SessionsInWorkerTotal)
	prometheus.MustRegister(repos.SessionsInWorkerCurrent)
	prometheus.MustRegister(repos.SessionsInWorkerUnitElapsed)
	prometheus.MustRegister(repos.SessionsInWorkerBatchSize)

	prometheus.MustRegister(middleware.DBIdleConnectionsGauge)
	prometheus.MustRegister(middleware.DBActiveConnectionsGauge)
	prometheus.MustRegister(middleware.DBOpenConnectionsGauge)
	prometheus.MustRegister(middleware.DBWaitCountGauge)
	prometheus.MustRegister(middleware.DBWaitSecondsGauge)
}
