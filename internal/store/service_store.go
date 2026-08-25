package store

import (
	"lightweightrpc/internal/model"
)

func (s *MemoryStore) CreateService(x *model.Service) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.services {
		if exist.Name == x.Name {
			return ErrConflict
		}
	}
	s.services[x.ID] = x
	return nil
}

func (s *MemoryStore) GetService(id string) (*model.Service, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	x, ok := s.services[id]
	if !ok {
		return nil, ErrNotFound
	}
	return x, nil
}

func (s *MemoryStore) GetServiceByName(v string) (*model.Service, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, x := range s.services {
		if x.Name == v {
			return x, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListServices() []*model.Service {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Service, 0, len(s.services))
	for _, x := range s.services {
		list = append(list, x)
	}
	return list
}

func (s *MemoryStore) UpdateService(x *model.Service) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.services[x.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.services {
		if exist.ID != x.ID && exist.Name == x.Name {
			return ErrConflict
		}
	}
	s.services[x.ID] = x
	return nil
}

func (s *MemoryStore) DeleteService(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.services[id]; !ok {
		return ErrNotFound
	}
	delete(s.services, id)
	return nil
}
