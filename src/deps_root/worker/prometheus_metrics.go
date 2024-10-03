package worker

import (
	"nosebook/src/lib/worker"

	"github.com/prometheus/client_golang/prometheus"
)

type prometheusMetrics struct {
	slug                string
	flushedBufferSize   *prometheus.HistogramVec
	elapsedSeconds      *prometheus.HistogramVec
	flushElapsedSeconds *prometheus.HistogramVec
}

var FlushedBufferSize = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "app_worker_buffer_flushed_size",
		Help:    "Buffer size flushed",
		Buckets: prometheus.ExponentialBuckets(1, 2, 10),
	},
	[]string{
		"slug",
	},
)

var SendElapsedSeconds = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "app_worker_buffer_send_elapsed_seconds",
		Help:    "Elapsed seconds buffer takes to complete one send",
		Buckets: prometheus.ExponentialBuckets(0.001, 2, 12),
	},
	[]string{
		"slug",
	},
)

var FlushElapsedSeconds = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "app_worker_buffer_flush_elapsed_seconds",
		Help:    "Elapsed seconds buffer flush implementation takes to complete",
		Buckets: prometheus.ExponentialBuckets(0.001, 2, 12),
	},
	[]string{
		"slug",
	},
)

func newPrometheusMetrics(slug string) *prometheusMetrics {
	return &prometheusMetrics{slug: slug}
}

func (this *prometheusMetrics) FlushedBufferSize(size float64) {
	FlushedBufferSize.WithLabelValues(this.slug).Observe(size)
}
func (this *prometheusMetrics) ElapsedSeconds(seconds float64) {
	SendElapsedSeconds.WithLabelValues(this.slug).Observe(seconds)
}
func (this *prometheusMetrics) ElapsedFlushSeconds(seconds float64) {
	FlushElapsedSeconds.WithLabelValues(this.slug).Observe(seconds)
}

type prometheusMetricsOpt struct {
	slug string
}

func (this *prometheusMetricsOpt) Metrics() worker.Metrics { return newPrometheusMetrics(this.slug) }
func (this *prometheusMetricsOpt) BufferSize() int         { return 0 }

func UsePrometheusMetrics(slug string) func() worker.BufferOpt {
	return func() worker.BufferOpt {
		return &prometheusMetricsOpt{
			slug: slug,
		}
	}
}
