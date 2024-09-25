package roothttp

import (
	"context"
	"log"
)

type ShutdownFn func()

func (this *RootHTTP) shutdown() {
	log.Println("Shutting down the server...")

	if err := this.traceProvider.Shutdown(context.Background()); err != nil {
		log.Printf("Error shutting down tracer provider: %v", err)
	}

	for _, shutdown := range this.shutdowns {
		shutdown()
	}
}
