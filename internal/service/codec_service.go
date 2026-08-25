package service

import (
	"sort"
	"time"

	"lightweightrpc/internal/model"
	"lightweightrpc/pkg/idgen"
)

func (s *Service) CreateCodec(input model.Codec) (*model.Codec, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetCodecByName(input.Name); err == nil {
		return nil, model.NewValidationError("name", "已存在")
	}
	now := time.Now()
	input.ID = idgen.Hex()
	input.CreatedAt = now
	input.UpdatedAt = now
	if err := s.store.CreateCodec(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *Service) GetCodec(id string) (*model.Codec, error) {
	return s.store.GetCodec(id)
}

func (s *Service) ListCodecs(filter model.CodecFilter, page, size int) ([]*model.Codec, int, error) {
	all := s.store.ListCodecs()
	matched := make([]*model.Codec, 0, len(all))
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
		return []*model.Codec{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateCodec(id string, input model.Codec) (*model.Codec, error) {
	exist, err := s.store.GetCodec(id)
	if err != nil {
		return nil, err
	}
	exist.Name = input.Name
	exist.Type = input.Type
	exist.Version = input.Version
	exist.Schema = input.Schema
	if err := exist.Validate(); err != nil {
		return nil, err
	}
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateCodec(exist); err != nil {
		return nil, err
	}
	return exist, nil
}

func (s *Service) DeleteCodec(id string) error {
	return s.store.DeleteCodec(id)
}

func (s *Service) TransitionCodec(id, target string) (*model.Codec, error) {
	if !model.CodecValidStatus(target) {
		return nil, model.NewValidationError("status", "目标状态不合法")
	}
	exist, err := s.store.GetCodec(id)
	if err != nil {
		return nil, err
	}
	if !model.CodecCanTransition(exist.Status, target) {
		return nil, model.NewValidationError("status", "不允许从 "+exist.Status+" 流转到 "+target)
	}
	exist.Status = target
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateCodec(exist); err != nil {
		return nil, err
	}
	return exist, nil
}
