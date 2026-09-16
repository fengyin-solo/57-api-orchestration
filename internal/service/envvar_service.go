package service

import (
	"sort"
	"time"

	"orchestr/internal/model"
	"orchestr/pkg/idgen"
)

func (s *Service) CreateEnvVar(input model.EnvVar) (*model.EnvVar, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetFlow(input.FlowID); err != nil {
		return nil, model.NewValidationError("flow_id", "流程不存在")
	}
	now := time.Now().Format(time.RFC3339)
	e := &model.EnvVar{
		ID:        idgen.Hex(),
		FlowID:    input.FlowID,
		Key:       input.Key,
		Value:     input.Value,
		IsSecret:  input.IsSecret,
		Status:    input.Status,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.store.CreateEnvVar(e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *Service) GetEnvVar(id string) (*model.EnvVar, error) {
	return s.store.GetEnvVar(id)
}

func (s *Service) ListEnvVars(filter model.EnvVarFilter, page, size int) ([]*model.EnvVar, int, error) {
	all := s.store.ListEnvVars()
	matched := make([]*model.EnvVar, 0, len(all))
	for _, e := range all {
		if filter.Match(e) {
			matched = append(matched, e)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt > matched[j].CreatedAt
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.EnvVar{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateEnvVar(id string, input model.EnvVar) (*model.EnvVar, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	e, err := s.store.GetEnvVar(id)
	if err != nil {
		return nil, err
	}
	e.FlowID = input.FlowID
	e.Key = input.Key
	e.Value = input.Value
	e.IsSecret = input.IsSecret
	e.Status = input.Status
	e.UpdatedAt = time.Now().Format(time.RFC3339)
	if err := s.store.UpdateEnvVar(e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *Service) DeleteEnvVar(id string) error {
	return s.store.DeleteEnvVar(id)
}

func (s *Service) BuildEnvMap(flowID string) map[string]string {
	vars := s.store.ListEnvVarsByFlowID(flowID)
	m := make(map[string]string)
	for _, v := range vars {
		if v.Status == model.EnvVarActive {
			m[v.Key] = v.Value
		}
	}
	return m
}
