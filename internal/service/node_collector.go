package service

import (
	"context"
	"sync"

	"lightweightrpc/internal/store"
)

func CollectNodeNames(ctx context.Context, names []string) ([]string, error) {
	out := make(chan string)
	errs := make(chan error)
	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		defer wg.Done()
		store.StreamNodeNames(names, out, errs)
	}()
	collected := make([]string, 0, len(names))
	for {
		select {
		case name := <-out:
			collected = append(collected, name)
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}
