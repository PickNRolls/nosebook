package worker

type Metrics interface {
	FlushedBufferSize(size float64)
	ElapsedSeconds(seconds float64)
	ElapsedFlushSeconds(seconds float64)
}

type noopMetrics struct{}

func (this *noopMetrics) FlushedBufferSize(size float64)      {}
func (this *noopMetrics) ElapsedSeconds(seconds float64)      {}
func (this *noopMetrics) ElapsedFlushSeconds(seconds float64) {}
