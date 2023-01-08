package workergroup

import (
	"context"
	"log"
	"sync"
)

type Worker func(ctx context.Context)

type WorkerGroup struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	wg        *sync.WaitGroup
}

func New() *WorkerGroup {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerGroup{
		ctx:       ctx,
		ctxCancel: cancel,
		wg:        &sync.WaitGroup{},
	}
}

func (wgx *WorkerGroup) Add(name string, fn Worker) {
	log.Printf("Adding worker %s", name)
	wgx.wg.Add(1)
	go func() {
		defer wgx.wg.Done()
		log.Printf("Worker: starting %s", name)
		fn(wgx.ctx)
		log.Printf("Worker: finished %s", name)
	}()
	log.Printf("Worker %s added", name)
}

func (wgx *WorkerGroup) Wait() {
	wgx.wg.Wait()
}

func (wgx *WorkerGroup) Stop() {
	wgx.ctxCancel()
}
