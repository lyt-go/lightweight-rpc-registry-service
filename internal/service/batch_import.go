package service

import (
	"sync"

	"lightweightrpc/internal/store"
)

type ImportAudit struct {
	mu      sync.Mutex
	success []string
}

func (a *ImportAudit) MarkSuccess(name string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.success = append(a.success, name)
}

func (a *ImportAudit) Count() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return len(a.success)
}

type BatchImporter struct {
	repo  *store.BatchRepository
	audit *ImportAudit
}

func NewBatchImporter(repo *store.BatchRepository, audit *ImportAudit) *BatchImporter {
	return &BatchImporter{repo: repo, audit: audit}
}

func (i *BatchImporter) Import(names []string) error {
	for _, name := range names {
		i.audit.MarkSuccess(name)
	}
	return i.repo.Persist(names)
}
