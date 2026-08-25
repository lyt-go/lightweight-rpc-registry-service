package store

import (
	"context"
	"sync"
	"time"
)

type ContextProbe struct {
	mu         sync.Mutex
	remembered context.Context
	calls      int
	delay      time.Duration
}

func NewContextProbe(delay time.Duration) *ContextProbe { return &ContextProbe{delay: delay} }

func (p *ContextProbe) Call(ctx context.Context) error {
	p.mu.Lock()
	p.calls++
	if p.remembered == nil {
		p.remembered = ctx
	}
	active := p.remembered
	p.mu.Unlock()
	select {
	case <-time.After(p.delay):
		return nil
	case <-active.Done():
		return active.Err()
	}
}

func (p *ContextProbe) Calls() int { p.mu.Lock(); defer p.mu.Unlock(); return p.calls }
