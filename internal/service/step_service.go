package service

import (
	"sort"
	"time"

	"orchestr/internal/model"
	"orchestr/pkg/idgen"
)

func (s *Service) CreateStep(input model.Step) (*model.Step, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetFlow(input.FlowID); err != nil {
		return nil, model.NewValidationError("flow_id", "流程不存在")
	}
	if _, err := s.store.GetService(input.ServiceID); err != nil {
		return nil, model.NewValidationError("service_id", "服务不存在")
	}
	now := time.Now().Format(time.RFC3339)
	st := &model.Step{
		ID:            idgen.Hex(),
		FlowID:        input.FlowID,
		ServiceID:     input.ServiceID,
		Method:        input.Method,
		Path:          input.Path,
		ParamsMapping: input.ParamsMapping,
		Order:         input.Order,
		TimeoutMs:     input.TimeoutMs,
		RetryCount:    input.RetryCount,
		Status:        input.Status,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := s.store.CreateStep(st); err != nil {
		return nil, err
	}
	_ = s.logAudit("system", "create_step", "step", st.ID, "创建步骤")
	return st, nil
}

func (s *Service) GetStep(id string) (*model.Step, error) {
	return s.store.GetStep(id)
}

func (s *Service) ListSteps(filter model.StepFilter, page, size int) ([]*model.Step, int, error) {
	all := s.store.ListSteps()
	matched := make([]*model.Step, 0, len(all))
	for _, st := range all {
		if filter.Match(st) {
			matched = append(matched, st)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].Order < matched[j].Order
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Step{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateStep(id string, input model.Step) (*model.Step, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	st, err := s.store.GetStep(id)
	if err != nil {
		return nil, err
	}
	st.FlowID = input.FlowID
	st.ServiceID = input.ServiceID
	st.Method = input.Method
	st.Path = input.Path
	st.ParamsMapping = input.ParamsMapping
	st.Order = input.Order
	st.TimeoutMs = input.TimeoutMs
	st.RetryCount = input.RetryCount
	st.Status = input.Status
	st.UpdatedAt = time.Now().Format(time.RFC3339)
	if err := s.store.UpdateStep(st); err != nil {
		return nil, err
	}
	_ = s.logAudit("system", "update_step", "step", st.ID, "更新步骤")
	return st, nil
}

func (s *Service) DeleteStep(id string) error {
	if err := s.store.DeleteStep(id); err != nil {
		return err
	}
	_ = s.logAudit("system", "delete_step", "step", id, "删除步骤")
	return nil
}

func (s *Service) BatchCreateSteps(inputs []model.Step) ([]*model.Step, error) {
	result := make([]*model.Step, 0, len(inputs))
	for _, input := range inputs {
		st, err := s.CreateStep(input)
		if err != nil {
			return nil, err
		}
		result = append(result, st)
	}
	return result, nil
}
