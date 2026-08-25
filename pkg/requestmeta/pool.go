package requestmeta

import "sync"

type Metadata struct {
	Tenant string
	Trace  string
}

type Pool struct {
	mu   sync.Mutex
	item *Metadata
}

func (p *Pool) Acquire() *Metadata {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.item == nil {
		return &Metadata{}
	}
	item := p.item
	p.item = nil
	return item
}

func (p *Pool) Release(item *Metadata) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.item = item
}
