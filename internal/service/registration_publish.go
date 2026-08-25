package service

import (
	"fmt"
	"sync"

	"lightweightrpc/internal/store"
)

type RegistrationPublisher interface{ Publish(string) error }

type SuccessAudit struct {
	mu    sync.Mutex
	count int
}

func (a *SuccessAudit) Mark()      { a.mu.Lock(); defer a.mu.Unlock(); a.count++ }
func (a *SuccessAudit) Count() int { a.mu.Lock(); defer a.mu.Unlock(); return a.count }

type RegistrationWorkflow struct {
	repo      *store.RegistrationRepository
	publisher RegistrationPublisher
	audit     *SuccessAudit
}

func NewRegistrationWorkflow(repo *store.RegistrationRepository, publisher RegistrationPublisher, audit *SuccessAudit) *RegistrationWorkflow {
	return &RegistrationWorkflow{repo: repo, publisher: publisher, audit: audit}
}

func (w *RegistrationWorkflow) Register(serviceID string) error {
	var err error
	for attempt := 1; attempt <= 2; attempt++ {
		eventID := fmt.Sprintf("%s:%d", serviceID, attempt)
		w.repo.CommitWithEvent(eventID)
		w.audit.Mark()
		err = w.publisher.Publish(eventID)
		if err == nil {
			return nil
		}
	}
	return err
}
