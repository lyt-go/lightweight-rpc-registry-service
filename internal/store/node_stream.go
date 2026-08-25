package store

import (
	"context"
	"fmt"
)

// StreamNodeNames 逐个将 names 中的节点名写入 out。
// 任一节点名为空即向 errs 发送错误并停止。
// 每次发送都监听 ctx：被取消时立即停止，避免在无人接收的无缓冲通道上永久阻塞。
// 返回前关闭 out 与 errs，使消费端能在正常结束或出错时及时收尾。
func StreamNodeNames(ctx context.Context, names []string, out chan<- string, errs chan<- error) {
	defer close(out)
	defer close(errs)
	for _, name := range names {
		if name == "" {
			select {
			case errs <- fmt.Errorf("node name is empty"):
			case <-ctx.Done():
			}
			return
		}
		select {
		case out <- name:
		case <-ctx.Done():
			return
		}
	}
}
