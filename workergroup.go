package workergroup

import (
	"context"
	"sync"

	"log/slog"
)

var Logger = slog.Default()

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
	Logger.Debug("Adding worker", "worker", name)
	wgx.wg.Add(1)
	go func() {
		defer wgx.wg.Done()
		Logger.Debug("Worker starting", "worker", name)
		fn(wgx.ctx)
		Logger.Debug("Worker finished", "worker", name)
	}()
	Logger.Debug("Worker added", "worker", name)
}

func (wgx *WorkerGroup) Wait() {
	wgx.wg.Wait()
}

func (wgx *WorkerGroup) Stop() {
	wgx.ctxCancel()
}
