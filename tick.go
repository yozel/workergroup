package workergroup

import (
	"context"
	"time"
)

func TickWorker(interval time.Duration, fn func()) Worker {
	tick := time.NewTicker(interval)
	immediate := make(chan struct{}, 1)
	immediate <- struct{}{}
	return func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
				fn()
			case <-immediate:
				fn()
			}
		}
	}
}
