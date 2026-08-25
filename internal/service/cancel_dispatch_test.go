package service

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"lightweightrpc/internal/worker"
)

func TestDispatchCancellationStopsRetriesBeforeShutdown(t *testing.T) {
	dispatcher := NewCancelDispatcher(&worker.RetryWorker{})
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Millisecond)
	defer cancel()
	var calls atomic.Int64
	dispatcher.Dispatch(ctx, func() error {
		calls.Add(1)
		return errors.New("node unavailable")
	})
	afterReturn := calls.Load()
	time.Sleep(12 * time.Millisecond)
	if got := calls.Load(); got != afterReturn {
		t.Fatalf("retry calls kept growing after request cancellation: before=%d after=%d", afterReturn, got)
	}
	if !dispatcher.Shutdown(20 * time.Millisecond) {
		t.Fatal("shutdown timed out waiting for a canceled request worker")
	}
}
