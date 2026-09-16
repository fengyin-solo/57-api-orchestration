package service

import (
	"sort"
	"time"

	"orchestr/internal/model"
	"orchestr/pkg/idgen"
)

func (s *Service) CreateDependency(input model.Dependency) (*model.Dependency, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetFlow(input.FlowID); err != nil {
		return nil, model.NewValidationError("flow_id", "流程不存在")
	}
	if _, err := s.store.GetStep(input.StepID); err != nil {
		return nil, model.NewValidationError("step_id", "步骤不存在")
	}
	if _, err := s.store.GetStep(input.DependsOnStepID); err != nil {
		return nil, model.NewValidationError("depends_on_step_id", "依赖步骤不存在")
	}
	now := time.Now().Format(time.RFC3339)
	d := &model.Dependency{
		ID:              idgen.Hex(),
		FlowID:          input.FlowID,
		StepID:          input.StepID,
		DependsOnStepID: input.DependsOnStepID,
		Condition:       input.Condition,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	if err := s.store.CreateDependency(d); err != nil {
		return nil, err
	}
	_ = s.logAudit("system", "create_dependency", "dependency", d.ID, "创建依赖")
	return d, nil
}

func (s *Service) GetDependency(id string) (*model.Dependency, error) {
	return s.store.GetDependency(id)
}

func (s *Service) ListDependencies(filter model.DependencyFilter, page, size int) ([]*model.Dependency, int, error) {
	all := s.store.ListDependencies()
	matched := make([]*model.Dependency, 0, len(all))
	for _, d := range all {
		if filter.Match(d) {
			matched = append(matched, d)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt > matched[j].CreatedAt
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Dependency{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateDependency(id string, input model.Dependency) (*model.Dependency, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	d, err := s.store.GetDependency(id)
	if err != nil {
		return nil, err
	}
	d.FlowID = input.FlowID
	d.StepID = input.StepID
	d.DependsOnStepID = input.DependsOnStepID
	d.Condition = input.Condition
	d.UpdatedAt = time.Now().Format(time.RFC3339)
	if err := s.store.UpdateDependency(d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Service) DeleteDependency(id string) error {
	return s.store.DeleteDependency(id)
}
