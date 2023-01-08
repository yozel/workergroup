package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/yozel/workergroup"
)

func main() {
	wgx := workergroup.New()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	go func() {
		<-sigs
		log.Println("shutting down")
		wgx.Stop()
	}()

	wgx.Add("worker1", func(ctx context.Context) {
		log.Println("worker 1 started")
		for i := 0; i < 10; i++ {
			select {
			case <-ctx.Done():
				return
			default:
				log.Println("worker1 running")
				time.Sleep(1 * time.Second)
			}
		}
		log.Println("worker1 stopped")
	})

	wgx.Add("worker2", func(ctx context.Context) {
		log.Println("worker 2 started")
		for i := 0; i < 5; i++ {
			select {
			case <-ctx.Done():
				return
			default:
				log.Println("worker2 running")
				time.Sleep(1 * time.Second)
			}
		}
		log.Println("worker2 stopped")
	})

	server := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("hello"))
		}),
	}

	wgx.Add("worker3", workergroup.HttpServerWorker(server))

	wgx.Wait()
}
