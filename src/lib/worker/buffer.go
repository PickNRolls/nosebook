package worker

type Buffer[S any, R any, W any] struct {
	doFlush <-chan W
	flush   Flush[S, R]
	send    chan send[S, R]
	done    <-chan struct{}
}

type send[S any, R any] struct {
	value    S
	receiver chan<- R
}
type Flush[S any, R any] func(values []S) R

func NewBuffer[S any, R any, W any](flush Flush[S, R], doFlush <-chan W, done <-chan struct{}, buffer int) *Buffer[S, R, W] {
	return &Buffer[S, R, W]{
		doFlush: doFlush,
		flush:   flush,
		send:    make(chan send[S, R], buffer),
		done:    done,
	}
}

func (this *Buffer[S, R, W]) Run() {
	for {
		select {
		case <-this.done:
			break

		case <-this.doFlush:
			values := []S{}
			receivers := []chan<- R{}
			for range len(this.send) {
				send := <-this.send
				values = append(values, send.value)
				receivers = append(receivers, send.receiver)
			}

			out := this.flush(values)

			for _, receiver := range receivers {
				receiver <- out
				close(receiver)
			}
		}
	}
}

func (this *Buffer[S, R, W]) Send(value S) R {
	ch := make(chan R)
	this.send <- send[S, R]{
		value:    value,
		receiver: ch,
	}
	return <-ch
}
