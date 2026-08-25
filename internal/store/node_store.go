package store

import (
	"lightweightrpc/internal/model"
)

func (s *MemoryStore) CreateNode(x *model.Node) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.nodes {
		if exist.Address == x.Address {
			return ErrConflict
		}
	}
	s.nodes[x.ID] = x
	return nil
}

func (s *MemoryStore) GetNode(id string) (*model.Node, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	x, ok := s.nodes[id]
	if !ok {
		return nil, ErrNotFound
	}
	return x, nil
}

func (s *MemoryStore) GetNodeByAddress(v string) (*model.Node, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, x := range s.nodes {
		if x.Address == v {
			return x, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListNodes() []*model.Node {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Node, 0, len(s.nodes))
	for _, x := range s.nodes {
		list = append(list, x)
	}
	return list
}

func (s *MemoryStore) UpdateNode(x *model.Node) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.nodes[x.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.nodes {
		if exist.ID != x.ID && exist.Address == x.Address {
			return ErrConflict
		}
	}
	s.nodes[x.ID] = x
	return nil
}

func (s *MemoryStore) DeleteNode(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.nodes[id]; !ok {
		return ErrNotFound
	}
	delete(s.nodes, id)
	return nil
}
