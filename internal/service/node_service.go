package service

import (
	"sort"
	"time"

	"lightweightrpc/internal/model"
	"lightweightrpc/pkg/idgen"
)

func (s *Service) CreateNode(input model.Node) (*model.Node, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetNodeByAddress(input.Address); err == nil {
		return nil, model.NewValidationError("address", "已存在")
	}
	now := time.Now()
	input.ID = idgen.Hex()
	input.CreatedAt = now
	input.UpdatedAt = now
	if err := s.store.CreateNode(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *Service) GetNode(id string) (*model.Node, error) {
	return s.store.GetNode(id)
}

func (s *Service) ListNodes(filter model.NodeFilter, page, size int) ([]*model.Node, int, error) {
	all := s.store.ListNodes()
	matched := make([]*model.Node, 0, len(all))
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
		return []*model.Node{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateNode(id string, input model.Node) (*model.Node, error) {
	exist, err := s.store.GetNode(id)
	if err != nil {
		return nil, err
	}
	exist.Name = input.Name
	exist.Address = input.Address
	exist.Region = input.Region
	exist.Weight = input.Weight
	if err := exist.Validate(); err != nil {
		return nil, err
	}
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateNode(exist); err != nil {
		return nil, err
	}
	return exist, nil
}

func (s *Service) DeleteNode(id string) error {
	return s.store.DeleteNode(id)
}

func (s *Service) TransitionNode(id, target string) (*model.Node, error) {
	if !model.NodeValidStatus(target) {
		return nil, model.NewValidationError("status", "目标状态不合法")
	}
	exist, err := s.store.GetNode(id)
	if err != nil {
		return nil, err
	}
	if !model.NodeCanTransition(exist.Status, target) {
		return nil, model.NewValidationError("status", "不允许从 "+exist.Status+" 流转到 "+target)
	}
	exist.Status = target
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateNode(exist); err != nil {
		return nil, err
	}
	return exist, nil
}
