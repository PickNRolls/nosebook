package worker

import (
	"nosebook/src/lib/clock"
)

type Buffer[S any, R any, W any] struct {
	metrics    Metrics
	flushEmpty bool
	doFlush    <-chan W
	flush      Flush[S, R]
	send       chan send[S, R]
	done       <-chan struct{}
}

type send[S any, R any] struct {
	value    S
	receiver chan<- R
}

type Flush[S any, R any] func(values []S) R

func NewBuffer[S any, R any, W any](flush Flush[S, R], doFlush <-chan W, done <-chan struct{}, bufferSize int, optFns ...func() BufferOpt) *Buffer[S, R, W] {
	flushEmpty := false
	var metrics Metrics = nil

	for _, optFn := range optFns {
		opt := optFn()

		if opt.FlushEmpty() {
			flushEmpty = true
		}

		if opt.Metrics() != nil && metrics == nil {
			metrics = opt.Metrics()
		}
	}

	if metrics == nil {
		metrics = &noopMetrics{}
	}

	metrics.Register()

	return &Buffer[S, R, W]{
		metrics:    metrics,
		flushEmpty: flushEmpty,
		doFlush:    doFlush,
		flush:      flush,
		send:       make(chan send[S, R], bufferSize),
		done:       done,
	}
}

func (this *Buffer[S, R, W]) Run() {
	for {
		select {
		case <-this.done:
			break

		case <-this.doFlush:
			if len(this.send) == 0 && !this.flushEmpty {
				continue
			}

			this.metrics.FlushedBufferSize(float64(len(this.send)))

			values := []S{}
			receivers := []chan<- R{}
			for range len(this.send) {
				send := <-this.send
				values = append(values, send.value)
				receivers = append(receivers, send.receiver)
			}

			beforeFlush := clock.Now()
			out := this.flush(values)
			this.metrics.ElapsedFlushSeconds(clock.Now().Sub(beforeFlush).Seconds())

			for _, receiver := range receivers {
				receiver <- out
				close(receiver)
			}
		}
	}
}

func (this *Buffer[S, R, W]) Send(value S) R {
	start := clock.Now()
	defer func() {
		this.metrics.ElapsedSeconds(clock.Now().Sub(start).Seconds())
	}()

	ch := make(chan R)
	this.send <- send[S, R]{
		value:    value,
		receiver: ch,
	}
	return <-ch
}
