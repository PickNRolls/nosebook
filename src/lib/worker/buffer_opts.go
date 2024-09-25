package worker

type BufferOpt interface {
	Metrics() Metrics
	BufferSize() int
	Done() <-chan struct{}
}

type bufferSizeOpt struct {
	size int
}

func (this *bufferSizeOpt) Metrics() Metrics      { return nil }
func (this *bufferSizeOpt) BufferSize() int       { return this.size }
func (this *bufferSizeOpt) Done() <-chan struct{} { return nil }

func BufferSize(size int) BufferOpt {
	return &bufferSizeOpt{
		size: size,
	}
}
