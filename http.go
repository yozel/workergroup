package workergroup

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"
)

func HttpServerWorker(server *http.Server) Worker {
	return startStopWorker(
		func() {
			if err := server.ListenAndServe(); err != nil {
				if err != http.ErrServerClosed {
					log.Printf("error while listening http server: %s", err)
				}
			}
		},
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := server.Shutdown(ctx); err != nil {
				log.Printf("error while stopping http server: %s", err)
			}
		},
	)
}

func startStopWorker(run func(), stop func()) Worker {
	return func(ctx context.Context) {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			run()
		}()
		<-ctx.Done()
		stop()
		wg.Wait()
	}
}
