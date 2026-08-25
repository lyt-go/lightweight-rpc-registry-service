package service

import (
	"sort"
	"time"

	"lightweightrpc/internal/model"
	"lightweightrpc/pkg/idgen"
)

func (s *Service) CreateRetryPolicy(input model.RetryPolicy) (*model.RetryPolicy, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetRetryPolicyByName(input.Name); err == nil {
		return nil, model.NewValidationError("name", "已存在")
	}
	now := time.Now()
	input.ID = idgen.Hex()
	input.CreatedAt = now
	input.UpdatedAt = now
	if err := s.store.CreateRetryPolicy(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *Service) GetRetryPolicy(id string) (*model.RetryPolicy, error) {
	return s.store.GetRetryPolicy(id)
}

func (s *Service) ListRetryPolicys(filter model.RetryPolicyFilter, page, size int) ([]*model.RetryPolicy, int, error) {
	all := s.store.ListRetryPolicys()
	matched := make([]*model.RetryPolicy, 0, len(all))
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
		return []*model.RetryPolicy{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateRetryPolicy(id string, input model.RetryPolicy) (*model.RetryPolicy, error) {
	exist, err := s.store.GetRetryPolicy(id)
	if err != nil {
		return nil, err
	}
	exist.Name = input.Name
	exist.MaxRetries = input.MaxRetries
	exist.BaseTimeoutMs = input.BaseTimeoutMs
	exist.Backoff = input.Backoff
	exist.Jitter = input.Jitter
	if err := exist.Validate(); err != nil {
		return nil, err
	}
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateRetryPolicy(exist); err != nil {
		return nil, err
	}
	return exist, nil
}

func (s *Service) DeleteRetryPolicy(id string) error {
	return s.store.DeleteRetryPolicy(id)
}

func (s *Service) TransitionRetryPolicy(id, target string) (*model.RetryPolicy, error) {
	if !model.RetryPolicyValidStatus(target) {
		return nil, model.NewValidationError("status", "目标状态不合法")
	}
	exist, err := s.store.GetRetryPolicy(id)
	if err != nil {
		return nil, err
	}
	if !model.RetryPolicyCanTransition(exist.Status, target) {
		return nil, model.NewValidationError("status", "不允许从 "+exist.Status+" 流转到 "+target)
	}
	exist.Status = target
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateRetryPolicy(exist); err != nil {
		return nil, err
	}
	return exist, nil
}
