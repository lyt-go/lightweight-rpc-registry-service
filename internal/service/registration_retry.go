package service

import "lightweightrpc/internal/store"

type RegistrationRetrier struct{ repo *store.AttemptRepository }

func NewRegistrationRetrier(repo *store.AttemptRepository) *RegistrationRetrier {
	return &RegistrationRetrier{repo: repo}
}

func (r *RegistrationRetrier) Register(operation func() error) error {
	var err error
	for attempt := 0; attempt < 2; attempt++ {
		err = r.repo.Run(operation)
		if err == nil || !store.Retryable(err) {
			return err
		}
	}
	return err
}
