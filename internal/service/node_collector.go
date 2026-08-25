package service

import (
	"context"
	"sync"

	"lightweightrpc/internal/store"
)

// CollectNodeNames 并发收集节点名，等待 StreamNodeNames 写完并关闭 out/errs。
// 正常结束（out 关闭）时返回已收集的名字；生产端出错时直接把错误交给调用方。
// 上下文取消时立即返回 ctx.Err()，并由生产端通过 select 退出，避免 goroutine 泄漏。
func CollectNodeNames(ctx context.Context, names []string) ([]string, error) {
	out := make(chan string)
	errs := make(chan error)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		store.StreamNodeNames(ctx, names, out, errs)
	}()

	collected := make([]string, 0, len(names))
	outOpen := true
	errsOpen := true
	var firstErr error

	for outOpen || errsOpen {
		select {
		case name, ok := <-out:
			if !ok {
				// out 已关闭：正常收尾，不再读
				outOpen = false
				continue
			}
			collected = append(collected, name)
		case err, ok := <-errs:
			if !ok {
				errsOpen = false
				continue
			}
			// 生产端错误直接交给调用方，记录首个错误后继续等流收尾
			if firstErr == nil {
				firstErr = err
			}
		case <-ctx.Done():
			// 取消：生产端 select 在同一 ctx 上，会自行退出并关闭通道
			wg.Wait()
			return nil, ctx.Err()
		}
	}

	// 等待生产 goroutine 确认退出
	wg.Wait()
	if firstErr != nil {
		return collected, firstErr
	}
	return collected, nil
}
