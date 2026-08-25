package service

import (
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

const maxPublishAttempts = 2

func (w *RegistrationWorkflow) Register(serviceID string) error {
	// 注册提交与发送成功事件相互独立：注册只提交一次，
	// 仅对发送进行重试，避免重试导致重复注册与重复审计。
	eventID := serviceID
	w.repo.CommitWithEvent(eventID)

	var err error
	for attempt := 1; attempt <= maxPublishAttempts; attempt++ {
		err = w.publisher.Publish(eventID)
		if err == nil {
			w.audit.Mark()
			return nil
		}
	}
	return err
}
