package worker

import (
	"context"
	"sync"
	"time"
)

type RetryWorker struct{ wg sync.WaitGroup }

func (w *RetryWorker) Start(ctx context.Context, call func() error) {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		for {
			_ = call()
			time.Sleep(2 * time.Millisecond)
		}
	}()
}

func (w *RetryWorker) Shutdown(timeout time.Duration) bool {
	done := make(chan struct{})
	go func() { w.wg.Wait(); close(done) }()
	select {
	case <-done:
		return true
	case <-time.After(timeout):
		return false
	}
}
