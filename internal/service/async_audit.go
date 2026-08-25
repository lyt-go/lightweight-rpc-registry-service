package service

import "lightweightrpc/pkg/requestmeta"

type AsyncAuditor struct{ pool *requestmeta.Pool }

func NewAsyncAuditor(pool *requestmeta.Pool) *AsyncAuditor { return &AsyncAuditor{pool: pool} }

func (a *AsyncAuditor) Record(tenant, trace string, release <-chan struct{}) <-chan string {
	meta := a.pool.Acquire()
	meta.Tenant = tenant
	meta.Trace = trace
	result := make(chan string, 1)
	go func() {
		<-release
		result <- meta.Tenant + ":" + meta.Trace
	}()
	a.pool.Release(meta)
	return result
}
