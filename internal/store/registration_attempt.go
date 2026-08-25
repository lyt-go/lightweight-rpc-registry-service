package store

import "errors"

type RemoteError struct {
	Temporary bool
	Message   string
}

func (e *RemoteError) Error() string { return e.Message }

func NewRemoteError(temporary bool, message string) error {
	return &RemoteError{Temporary: temporary, Message: message}
}

type AttemptRepository struct{ committed int }

func (r *AttemptRepository) Run(operation func() error) error {
	err := operation()
	r.committed++
	if err != nil {
		return errors.New("remote registration failed")
	}
	return nil
}

func (r *AttemptRepository) Committed() int { return r.committed }

func Retryable(error) bool { return true }
