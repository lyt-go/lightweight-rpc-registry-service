package service

import (
	"errors"
	"sync"
	"testing"

	"lightweightrpc/internal/store"
)

type failOncePublisher struct {
	mu        sync.Mutex
	calls     int
	published map[string]struct{}
}

func (p *failOncePublisher) Publish(id string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.calls++
	if p.calls == 1 {
		return errors.New("publisher unavailable")
	}
	if p.published == nil {
		p.published = make(map[string]struct{})
	}
	p.published[id] = struct{}{}
	return nil
}

func (p *failOncePublisher) UniqueEvents() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.published)
}

func TestPublisherRetryKeepsSingleCommitAndAudit(t *testing.T) {
	repo := &store.RegistrationRepository{}
	publisher := &failOncePublisher{}
	audit := &SuccessAudit{}
	workflow := NewRegistrationWorkflow(repo, publisher, audit)
	if err := workflow.Register("profile-api"); err != nil {
		t.Fatalf("registration should recover after publisher retry: %v", err)
	}
	if got := repo.Commits(); got != 1 {
		t.Fatalf("publisher retry repeated the registration commit: got %d want 1", got)
	}
	if got := audit.Count(); got != 1 {
		t.Fatalf("success audit was published %d times, want 1", got)
	}
	if got := publisher.UniqueEvents(); got != 1 {
		t.Fatalf("publisher observed %d event identities, want 1", got)
	}
}
