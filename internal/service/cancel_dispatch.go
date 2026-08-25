package service

import (
	"context"
	"time"

	"lightweightrpc/internal/worker"
)

type CancelDispatcher struct{ worker *worker.RetryWorker }

func NewCancelDispatcher(w *worker.RetryWorker) *CancelDispatcher {
	return &CancelDispatcher{worker: w}
}

func (d *CancelDispatcher) Dispatch(ctx context.Context, call func() error) {
	d.worker.Start(context.Background(), call)
	<-ctx.Done()
}

func (d *CancelDispatcher) Shutdown(timeout time.Duration) bool { return d.worker.Shutdown(timeout) }
