package service

import (
	"fmt"
	"sort"

	"orchestr/internal/model"
)

func (s *Service) CreateStepExecution(input model.StepExecution) (*model.StepExecution, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetExecution(input.ExecutionID); err != nil {
		return nil, model.NewValidationError("execution_id", "执行实例不存在")
	}
	if _, err := s.store.GetStep(input.StepID); err != nil {
		return nil, model.NewValidationError("step_id", "步骤不存在")
	}
	if err := s.store.CreateStepExecution(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *Service) GetStepExecution(id string) (*model.StepExecution, error) {
	return s.store.GetStepExecution(id)
}

func (s *Service) ListStepExecutions(filter model.StepExecutionFilter, page, size int) ([]*model.StepExecution, int, error) {
	all := s.store.ListStepExecutions()
	matched := make([]*model.StepExecution, 0, len(all))
	for _, se := range all {
		if filter.Match(se) {
			matched = append(matched, se)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].StartedAt > matched[j].StartedAt
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.StepExecution{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateStepExecution(id string, input model.StepExecution) (*model.StepExecution, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	se, err := s.store.GetStepExecution(id)
	if err != nil {
		return nil, err
	}
	se.Status = input.Status
	se.Response = input.Response
	se.DurationMs = input.DurationMs
	se.StartedAt = input.StartedAt
	se.FinishedAt = input.FinishedAt
	if err := s.store.UpdateStepExecution(se); err != nil {
		return nil, err
	}
	return se, nil
}

func (s *Service) DeleteStepExecution(id string) error {
	return s.store.DeleteStepExecution(id)
}

func (s *Service) TransitionStepExecution(id string, toStatus string) (*model.StepExecution, error) {
	se, err := s.store.GetStepExecution(id)
	if err != nil {
		return nil, err
	}
	if !model.StepExecutionCanTransition(se.Status, toStatus) {
		return nil, model.NewValidationError("status", fmt.Sprintf("不能从 %s 转为 %s", se.Status, toStatus))
	}
	se.Status = toStatus
	if err := s.store.UpdateStepExecution(se); err != nil {
		return nil, err
	}
	return se, nil
}
