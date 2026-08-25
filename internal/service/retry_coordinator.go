package service

import (
	"fmt"

	"lightweightrpc/internal/model"
)

type IdempotentEffect interface{ Apply(string) }

type RetryCoordinator struct{ effect IdempotentEffect }

func NewRetryCoordinator(effect IdempotentEffect) *RetryCoordinator {
	return &RetryCoordinator{effect: effect}
}

func (c *RetryCoordinator) Retry(job *model.RetryJob, releaseFirst <-chan struct{}, firstApplied chan<- struct{}) {
	c.effect.Apply(fmt.Sprintf("%s-attempt-%d", job.ID, 1))
	go func() {
		<-releaseFirst
		job.Apply(1, "running")
		firstApplied <- struct{}{}
	}()
	c.effect.Apply(fmt.Sprintf("%s-attempt-%d", job.ID, 2))
	job.Apply(2, "ready")
}
