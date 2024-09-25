package roothttp

import (
	"context"
	"log"
)

type ShutdownFn func()

type Runnable interface {
  Run()
  OnDone()
}

func (this *RootHTTP) shutdownRun(runnable Runnable) {
  go runnable.Run()
  this.shutdowns = append(this.shutdowns, runnable.OnDone)
}

func (this *RootHTTP) shutdown() {
	log.Println("Shutting down the server...")

	if err := this.traceProvider.Shutdown(context.Background()); err != nil {
		log.Printf("Error shutting down tracer provider: %v", err)
	}

	for _, shutdown := range this.shutdowns {
		shutdown()
	}
}
