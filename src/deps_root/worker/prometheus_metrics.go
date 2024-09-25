package worker

import (
	"nosebook/src/lib/worker"

	"github.com/prometheus/client_golang/prometheus"
)

type prometheusMetrics struct {
	flushedBufferSize   prometheus.Histogram
	elapsedSeconds      prometheus.Histogram
	flushElapsedSeconds prometheus.Histogram
}

func newPrometheusMetrics(slug string) *prometheusMetrics {
	out := &prometheusMetrics{
		flushedBufferSize: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Name:    "app_" + slug + "_flushed_buffer_size",
				Help:    "Buffer size flushed",
				Buckets: prometheus.ExponentialBuckets(1, 2, 10),
			},
		),

		elapsedSeconds: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Name:    "app_" + slug + "_buffer_elapsed_seconds",
				Help:    "Elapsed seconds buffer takes to complete one send",
				Buckets: prometheus.ExponentialBuckets(0.001, 2, 12),
			},
		),

		flushElapsedSeconds: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Name:    "app_" + slug + "_buffer_flush_elapsed_seconds",
				Help:    "Elapsed seconds buffer flush implementation takes to complete",
				Buckets: prometheus.ExponentialBuckets(0.001, 2, 12),
			},
		),
	}

	return out
}

func (this *prometheusMetrics) FlushedBufferSize(size float64) {
	this.flushedBufferSize.Observe(size)
}
func (this *prometheusMetrics) ElapsedSeconds(seconds float64) {
	this.elapsedSeconds.Observe(seconds)
}
func (this *prometheusMetrics) ElapsedFlushSeconds(seconds float64) {
	this.flushElapsedSeconds.Observe(seconds)
}
func (this *prometheusMetrics) Register() {
	prometheus.MustRegister(this.elapsedSeconds)
	prometheus.MustRegister(this.flushedBufferSize)
	prometheus.MustRegister(this.flushElapsedSeconds)
}

type prometheusMetricsOpt struct {
	slug string
}

func (this *prometheusMetricsOpt) FlushEmpty() bool { return false }
func (this *prometheusMetricsOpt) Metrics() worker.Metrics {
	return newPrometheusMetrics(this.slug)
}

func UsePrometheusMetrics(slug string) func() worker.BufferOpt {
	return func() worker.BufferOpt {
		return &prometheusMetricsOpt{
			slug: slug,
		}
	}
}
