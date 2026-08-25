package service

import (
	"context"

	"lightweightrpc/internal/store"
)

type ContextResolver struct{ probe *store.ContextProbe }

func NewContextResolver(probe *store.ContextProbe) *ContextResolver {
	return &ContextResolver{probe: probe}
}

func (r *ContextResolver) Resolve(ctx context.Context) error {
	var err error
	for attempt := 0; attempt < 3; attempt++ {
		err = r.probe.Call(ctx)
		if err == nil {
			return nil
		}
		// 上下文已被取消/超时：该次探测已立即失败，重试也不可能恢复，
		// 立即停住，避免连续空转三次并让旧的超时状态被带入后续请求。
		if ctx.Err() != nil {
			return err
		}
	}
	return err
}
