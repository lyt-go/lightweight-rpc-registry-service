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
			// 请求结束或服务关闭时，立即停止重试，不再发起下一次调用。
			select {
			case <-ctx.Done():
				return
			default:
			}
			_ = call()
			// 退避等待同样可取消：ctx 结束后无需等满 2ms 即可退出。
			select {
			case <-ctx.Done():
				return
			case <-time.After(2 * time.Millisecond):
			}
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
