package store

import (
	"lightweightrpc/internal/model"
)

func (s *MemoryStore) CreateRetryPolicy(x *model.RetryPolicy) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.retryPolicies {
		if exist.Name == x.Name {
			return ErrConflict
		}
	}
	s.retryPolicies[x.ID] = x
	return nil
}

func (s *MemoryStore) GetRetryPolicy(id string) (*model.RetryPolicy, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	x, ok := s.retryPolicies[id]
	if !ok {
		return nil, ErrNotFound
	}
	return x, nil
}

func (s *MemoryStore) GetRetryPolicyByName(v string) (*model.RetryPolicy, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, x := range s.retryPolicies {
		if x.Name == v {
			return x, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListRetryPolicys() []*model.RetryPolicy {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.RetryPolicy, 0, len(s.retryPolicies))
	for _, x := range s.retryPolicies {
		list = append(list, x)
	}
	return list
}

func (s *MemoryStore) UpdateRetryPolicy(x *model.RetryPolicy) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.retryPolicies[x.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.retryPolicies {
		if exist.ID != x.ID && exist.Name == x.Name {
			return ErrConflict
		}
	}
	s.retryPolicies[x.ID] = x
	return nil
}

func (s *MemoryStore) DeleteRetryPolicy(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.retryPolicies[id]; !ok {
		return ErrNotFound
	}
	delete(s.retryPolicies, id)
	return nil
}
