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
	// 把请求 ctx 透传给 worker：请求结束（ctx 取消）后重试必须停下，
	// 不能再用 context.Background() 让它无限重试下去。
	d.worker.Start(ctx, call)
	<-ctx.Done()
}

func (d *CancelDispatcher) Shutdown(timeout time.Duration) bool { return d.worker.Shutdown(timeout) }
