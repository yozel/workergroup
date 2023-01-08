package workergroup

import "context"

func ChannelWorker[T any](ch <-chan T, fn func(T)) Worker {
	return func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-ch:
				fn(msg)
			}
		}
	}
}
