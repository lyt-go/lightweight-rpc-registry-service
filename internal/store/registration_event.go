package store

import "sync"

type RegistrationRepository struct {
	mu      sync.Mutex
	commits int
	events  []string
}

func (r *RegistrationRepository) CommitWithEvent(eventID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.commits++
	r.events = append(r.events, eventID)
}

func (r *RegistrationRepository) Commits() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.commits
}
