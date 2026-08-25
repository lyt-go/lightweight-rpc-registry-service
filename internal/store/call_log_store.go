package store

import (
	"lightweightrpc/internal/model"
)

func (s *MemoryStore) CreateCallLog(x *model.CallLog) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.callLogs[x.ID] = x
	return nil
}

func (s *MemoryStore) GetCallLog(id string) (*model.CallLog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	x, ok := s.callLogs[id]
	if !ok {
		return nil, ErrNotFound
	}
	return x, nil
}

func (s *MemoryStore) ListCallLogs() []*model.CallLog {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.CallLog, 0, len(s.callLogs))
	for _, x := range s.callLogs {
		list = append(list, x)
	}
	return list
}

func (s *MemoryStore) DeleteCallLog(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.callLogs[id]; !ok {
		return ErrNotFound
	}
	delete(s.callLogs, id)
	return nil
}
