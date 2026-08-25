package service

import (
	"context"

	"lightweightrpc/internal/store"
)

type ContextResolver struct{ probe *store.ContextProbe }

func NewContextResolver(probe *store.ContextProbe) *ContextResolver {
	return &ContextResolver{probe: probe}
}

func (r *ContextResolver) Resolve(ctx context.Context) error {
	var err error
	for attempt := 0; attempt < 3; attempt++ {
		err = r.probe.Call(ctx)
		if err == nil {
			return nil
		}
	}
	return err
}
