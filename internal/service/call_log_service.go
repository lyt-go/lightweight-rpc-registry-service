package service

import (
	"sort"
	"time"

	"lightweightrpc/internal/model"
	"lightweightrpc/pkg/idgen"
)

func (s *Service) CreateCallLog(input model.CallLog) (*model.CallLog, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	input.ID = idgen.Hex()
	input.CreatedAt = now
	input.UpdatedAt = now
	if err := s.store.CreateCallLog(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *Service) GetCallLog(id string) (*model.CallLog, error) {
	return s.store.GetCallLog(id)
}

func (s *Service) ListCallLogs(filter model.CallLogFilter, page, size int) ([]*model.CallLog, int, error) {
	all := s.store.ListCallLogs()
	matched := make([]*model.CallLog, 0, len(all))
	for _, x := range all {
		if filter.Match(x) {
			matched = append(matched, x)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.CallLog{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) DeleteCallLog(id string) error {
	return s.store.DeleteCallLog(id)
}
