package worker

type BufferOpt interface {
	Metrics() Metrics
	BufferSize() int
}

type bufferSizeOpt struct {
	size int
}

func (this *bufferSizeOpt) Metrics() Metrics { return nil }
func (this *bufferSizeOpt) BufferSize() int  { return this.size }

func BufferSize(size int) BufferOpt {
	return &bufferSizeOpt{
		size: size,
	}
}
