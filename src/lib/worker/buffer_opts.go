package worker

type BufferOpt interface {
	FlushEmpty() bool
	Metrics() Metrics
}

type flushEmptyOpt struct{}

func (this *flushEmptyOpt) FlushEmpty() bool { return true }
func (this *flushEmptyOpt) Metrics() Metrics { return nil }

func FlushEmpty() BufferOpt {
	return &flushEmptyOpt{}
}
