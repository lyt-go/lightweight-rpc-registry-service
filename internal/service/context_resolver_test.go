package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"lightweightrpc/internal/store"
)

func TestCanceledLookupDoesNotPoisonNextRequestOrRetry(t *testing.T) {
	probe := store.NewContextProbe(30 * time.Millisecond)
	resolver := NewContextResolver(probe)
	first, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()
	if err := resolver.Resolve(first); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("first lookup should stop at its deadline: %v", err)
	}
	if got := probe.Calls(); got != 1 {
		t.Fatalf("canceled lookup kept retrying: got %d calls want 1", got)
	}
	if err := resolver.Resolve(context.Background()); err != nil {
		t.Fatalf("fresh lookup inherited the previous deadline: %v", err)
	}
}
