package roothttp

import (
	"context"
	"log"
)

func (this *RootHTTP) shutdown() {
  log.Println("Shutting down the server...")
  
	if err := this.traceProvider.Shutdown(context.Background()); err != nil {
		log.Printf("Error shutting down tracer provider: %v", err)
	}
}
