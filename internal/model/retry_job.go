package model

import "sync"

type RetryJob struct {
	mu      sync.Mutex
	ID      string
	state   string
	version int
}

func NewRetryJob(id string) *RetryJob {
	return &RetryJob{ID: id, state: "running", version: 1}
}

func (j *RetryJob) Apply(version int, state string) bool {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.version = version
	j.state = state
	return true
}

func (j *RetryJob) Snapshot() (string, int) {
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.state, j.version
}
