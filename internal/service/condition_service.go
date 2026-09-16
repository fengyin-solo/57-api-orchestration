package service

import (
	"sort"
	"time"

	"orchestr/internal/model"
	"orchestr/pkg/idgen"
)

func (s *Service) CreateCondition(input model.Condition) (*model.Condition, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetFlow(input.FlowID); err != nil {
		return nil, model.NewValidationError("flow_id", "流程不存在")
	}
	if _, err := s.store.GetStep(input.StepID); err != nil {
		return nil, model.NewValidationError("step_id", "步骤不存在")
	}
	now := time.Now().Format(time.RFC3339)
	c := &model.Condition{
		ID:          idgen.Hex(),
		FlowID:      input.FlowID,
		StepID:      input.StepID,
		Expression:  input.Expression,
		TrueStepID:  input.TrueStepID,
		FalseStepID: input.FalseStepID,
		Status:      input.Status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.store.CreateCondition(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) GetCondition(id string) (*model.Condition, error) {
	return s.store.GetCondition(id)
}

func (s *Service) ListConditions(filter model.ConditionFilter, page, size int) ([]*model.Condition, int, error) {
	all := s.store.ListConditions()
	matched := make([]*model.Condition, 0, len(all))
	for _, c := range all {
		if filter.Match(c) {
			matched = append(matched, c)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt > matched[j].CreatedAt
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Condition{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateCondition(id string, input model.Condition) (*model.Condition, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	c, err := s.store.GetCondition(id)
	if err != nil {
		return nil, err
	}
	c.FlowID = input.FlowID
	c.StepID = input.StepID
	c.Expression = input.Expression
	c.TrueStepID = input.TrueStepID
	c.FalseStepID = input.FalseStepID
	c.Status = input.Status
	c.UpdatedAt = time.Now().Format(time.RFC3339)
	if err := s.store.UpdateCondition(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) DeleteCondition(id string) error {
	return s.store.DeleteCondition(id)
}
