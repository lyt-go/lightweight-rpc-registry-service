package store

import (
	"lightweightrpc/internal/model"
)

func (s *MemoryStore) CreateMethod(x *model.Method) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.methods {
		if exist.Name == x.Name {
			return ErrConflict
		}
	}
	s.methods[x.ID] = x
	return nil
}

func (s *MemoryStore) GetMethod(id string) (*model.Method, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	x, ok := s.methods[id]
	if !ok {
		return nil, ErrNotFound
	}
	return x, nil
}

func (s *MemoryStore) GetMethodByName(v string) (*model.Method, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, x := range s.methods {
		if x.Name == v {
			return x, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListMethods() []*model.Method {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Method, 0, len(s.methods))
	for _, x := range s.methods {
		list = append(list, x)
	}
	return list
}

func (s *MemoryStore) UpdateMethod(x *model.Method) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.methods[x.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.methods {
		if exist.ID != x.ID && exist.Name == x.Name {
			return ErrConflict
		}
	}
	s.methods[x.ID] = x
	return nil
}

func (s *MemoryStore) DeleteMethod(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.methods[id]; !ok {
		return ErrNotFound
	}
	delete(s.methods, id)
	return nil
}
