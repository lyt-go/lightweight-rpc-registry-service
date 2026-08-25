package service

import (
	"sort"
	"time"

	"lightweightrpc/internal/model"
	"lightweightrpc/pkg/idgen"
)

func (s *Service) CreateMethod(input model.Method) (*model.Method, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetMethodByName(input.Name); err == nil {
		return nil, model.NewValidationError("name", "已存在")
	}
	now := time.Now()
	input.ID = idgen.Hex()
	input.CreatedAt = now
	input.UpdatedAt = now
	if err := s.store.CreateMethod(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *Service) GetMethod(id string) (*model.Method, error) {
	return s.store.GetMethod(id)
}

func (s *Service) ListMethods(filter model.MethodFilter, page, size int) ([]*model.Method, int, error) {
	all := s.store.ListMethods()
	matched := make([]*model.Method, 0, len(all))
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
		return []*model.Method{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateMethod(id string, input model.Method) (*model.Method, error) {
	exist, err := s.store.GetMethod(id)
	if err != nil {
		return nil, err
	}
	exist.ServiceID = input.ServiceID
	exist.Name = input.Name
	exist.InputType = input.InputType
	exist.OutputType = input.OutputType
	exist.TimeoutMs = input.TimeoutMs
	if err := exist.Validate(); err != nil {
		return nil, err
	}
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateMethod(exist); err != nil {
		return nil, err
	}
	return exist, nil
}

func (s *Service) DeleteMethod(id string) error {
	return s.store.DeleteMethod(id)
}

func (s *Service) TransitionMethod(id, target string) (*model.Method, error) {
	if !model.MethodValidStatus(target) {
		return nil, model.NewValidationError("status", "目标状态不合法")
	}
	exist, err := s.store.GetMethod(id)
	if err != nil {
		return nil, err
	}
	if !model.MethodCanTransition(exist.Status, target) {
		return nil, model.NewValidationError("status", "不允许从 "+exist.Status+" 流转到 "+target)
	}
	exist.Status = target
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateMethod(exist); err != nil {
		return nil, err
	}
	return exist, nil
}
