package workergroup

import (
	"context"
	"net/http"
	"sync"
	"time"
)

func HttpServerWorker(server *http.Server) Worker {
	return startStopWorker(
		func() {
			if err := server.ListenAndServe(); err != nil {
				if err != http.ErrServerClosed {
					Logger.Error("error while listening http server", "err", err)
				}
			}
		},
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := server.Shutdown(ctx); err != nil {
				Logger.Error("error while stopping http server", "err", err)
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
