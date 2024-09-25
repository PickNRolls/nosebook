package worker

import (
	"nosebook/src/lib/clock"
)

type Buffer[S any, R any] struct {
	metrics Metrics
	flush   Flush[S, R]
	send    chan send[S, R]
	done    <-chan struct{}
}

type send[S any, R any] struct {
	value    S
	receiver chan<- R
}

type Flush[S any, R any] func(values []S) R

func NewBuffer[S any, R any](flush Flush[S, R], optFns ...func() BufferOpt) *Buffer[S, R] {
	var metrics Metrics = nil
	var done <-chan struct{}
	bufferSize := 0

	for _, optFn := range optFns {
		opt := optFn()

		if opt.Metrics() != nil && metrics == nil {
			metrics = opt.Metrics()
		}

		ch := opt.Done()
		if ch != nil {
			done = ch
		}

		if opt.BufferSize() != 0 {
			bufferSize = opt.BufferSize()
		}
	}

	if bufferSize == 0 {
		bufferSize = 256
	}

	if metrics == nil {
		metrics = &noopMetrics{}
	}

	metrics.Register()

	return &Buffer[S, R]{
		metrics: metrics,
		flush:   flush,
		send:    make(chan send[S, R], bufferSize),
		done:    done,
	}
}

func (this *Buffer[S, R]) Run() {
	for {
		select {
		case <-this.done:
			break

		case value := <-this.send:
			this.metrics.FlushedBufferSize(float64(len(this.send)))

			values := []S{value.value}
			receivers := []chan<- R{value.receiver}
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

func (this *Buffer[S, R]) Send(value S) R {
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
