package store

import (
	"lightweightrpc/internal/model"
)

func (s *MemoryStore) CreateCodec(x *model.Codec) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.codecs {
		if exist.Name == x.Name {
			return ErrConflict
		}
	}
	s.codecs[x.ID] = x
	return nil
}

func (s *MemoryStore) GetCodec(id string) (*model.Codec, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	x, ok := s.codecs[id]
	if !ok {
		return nil, ErrNotFound
	}
	return x, nil
}

func (s *MemoryStore) GetCodecByName(v string) (*model.Codec, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, x := range s.codecs {
		if x.Name == v {
			return x, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListCodecs() []*model.Codec {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Codec, 0, len(s.codecs))
	for _, x := range s.codecs {
		list = append(list, x)
	}
	return list
}

func (s *MemoryStore) UpdateCodec(x *model.Codec) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.codecs[x.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.codecs {
		if exist.ID != x.ID && exist.Name == x.Name {
			return ErrConflict
		}
	}
	s.codecs[x.ID] = x
	return nil
}

func (s *MemoryStore) DeleteCodec(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.codecs[id]; !ok {
		return ErrNotFound
	}
	delete(s.codecs, id)
	return nil
}
