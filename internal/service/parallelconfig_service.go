package service

import (
	"sort"
	"time"

	"orchestr/internal/model"
	"orchestr/pkg/idgen"
)

func (s *Service) CreateParallelConfig(input model.ParallelConfig) (*model.ParallelConfig, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetFlow(input.FlowID); err != nil {
		return nil, model.NewValidationError("flow_id", "流程不存在")
	}
	for _, sid := range input.StepIDs {
		if _, err := s.store.GetStep(sid); err != nil {
			return nil, model.NewValidationError("step_ids", "步骤不存在: "+sid)
		}
	}
	now := time.Now().Format(time.RFC3339)
	p := &model.ParallelConfig{
		ID:        idgen.Hex(),
		FlowID:    input.FlowID,
		StepIDs:   input.StepIDs,
		Status:    input.Status,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.store.CreateParallelConfig(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) GetParallelConfig(id string) (*model.ParallelConfig, error) {
	return s.store.GetParallelConfig(id)
}

func (s *Service) ListParallelConfigs(filter model.ParallelConfigFilter, page, size int) ([]*model.ParallelConfig, int, error) {
	all := s.store.ListParallelConfigs()
	matched := make([]*model.ParallelConfig, 0, len(all))
	for _, p := range all {
		if filter.Match(p) {
			matched = append(matched, p)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt > matched[j].CreatedAt
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.ParallelConfig{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateParallelConfig(id string, input model.ParallelConfig) (*model.ParallelConfig, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	p, err := s.store.GetParallelConfig(id)
	if err != nil {
		return nil, err
	}
	p.FlowID = input.FlowID
	p.StepIDs = input.StepIDs
	p.Status = input.Status
	p.UpdatedAt = time.Now().Format(time.RFC3339)
	if err := s.store.UpdateParallelConfig(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) DeleteParallelConfig(id string) error {
	return s.store.DeleteParallelConfig(id)
}
