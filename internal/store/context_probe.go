package store

import (
	"context"
	"sync"
	"time"
)

type ContextProbe struct {
	mu    sync.Mutex
	calls int
	delay time.Duration
}

func NewContextProbe(delay time.Duration) *ContextProbe { return &ContextProbe{delay: delay} }

// Call 仅监听本次调用传入的 ctx，不跨请求复用任何已取消的上下文，
// 避免上一轮的超时状态泄漏到后续新请求中。
func (p *ContextProbe) Call(ctx context.Context) error {
	p.mu.Lock()
	p.calls++
	p.mu.Unlock()
	select {
	case <-time.After(p.delay):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *ContextProbe) Calls() int { p.mu.Lock(); defer p.mu.Unlock(); return p.calls }
