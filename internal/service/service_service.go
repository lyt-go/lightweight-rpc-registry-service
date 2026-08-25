package service

import (
	"sort"
	"time"

	"lightweightrpc/internal/model"
	"lightweightrpc/pkg/idgen"
)

func (s *Service) CreateService(input model.Service) (*model.Service, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetServiceByName(input.Name); err == nil {
		return nil, model.NewValidationError("name", "已存在")
	}
	now := time.Now()
	input.ID = idgen.Hex()
	input.CreatedAt = now
	input.UpdatedAt = now
	if err := s.store.CreateService(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *Service) GetService(id string) (*model.Service, error) {
	return s.store.GetService(id)
}

func (s *Service) ListServices(filter model.ServiceFilter, page, size int) ([]*model.Service, int, error) {
	all := s.store.ListServices()
	matched := make([]*model.Service, 0, len(all))
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
		return []*model.Service{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateService(id string, input model.Service) (*model.Service, error) {
	exist, err := s.store.GetService(id)
	if err != nil {
		return nil, err
	}
	exist.Name = input.Name
	exist.Version = input.Version
	exist.Host = input.Host
	exist.Port = input.Port
	exist.Protocol = input.Protocol
	if err := exist.Validate(); err != nil {
		return nil, err
	}
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateService(exist); err != nil {
		return nil, err
	}
	return exist, nil
}

func (s *Service) DeleteService(id string) error {
	return s.store.DeleteService(id)
}

func (s *Service) TransitionService(id, target string) (*model.Service, error) {
	if !model.ServiceValidStatus(target) {
		return nil, model.NewValidationError("status", "目标状态不合法")
	}
	exist, err := s.store.GetService(id)
	if err != nil {
		return nil, err
	}
	if !model.ServiceCanTransition(exist.Status, target) {
		return nil, model.NewValidationError("status", "不允许从 "+exist.Status+" 流转到 "+target)
	}
	exist.Status = target
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateService(exist); err != nil {
		return nil, err
	}
	return exist, nil
}
