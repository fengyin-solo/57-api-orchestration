package service

import (
	"sort"
	"time"

	"orchestr/internal/model"
	"orchestr/pkg/idgen"
)

func (s *Service) CreateRetryPolicy(input model.RetryPolicy) (*model.RetryPolicy, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now().Format(time.RFC3339)
	r := &model.RetryPolicy{
		ID:             idgen.Hex(),
		Name:           input.Name,
		MaxRetries:     input.MaxRetries,
		BackoffType:    input.BackoffType,
		InitialDelayMs: input.InitialDelayMs,
		Status:         input.Status,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	if err := s.store.CreateRetryPolicy(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) GetRetryPolicy(id string) (*model.RetryPolicy, error) {
	return s.store.GetRetryPolicy(id)
}

func (s *Service) ListRetryPolicies(filter model.RetryPolicyFilter, page, size int) ([]*model.RetryPolicy, int, error) {
	all := s.store.ListRetryPolicies()
	matched := make([]*model.RetryPolicy, 0, len(all))
	for _, r := range all {
		if filter.Match(r) {
			matched = append(matched, r)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt > matched[j].CreatedAt
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.RetryPolicy{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateRetryPolicy(id string, input model.RetryPolicy) (*model.RetryPolicy, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	r, err := s.store.GetRetryPolicy(id)
	if err != nil {
		return nil, err
	}
	r.Name = input.Name
	r.MaxRetries = input.MaxRetries
	r.BackoffType = input.BackoffType
	r.InitialDelayMs = input.InitialDelayMs
	r.Status = input.Status
	r.UpdatedAt = time.Now().Format(time.RFC3339)
	if err := s.store.UpdateRetryPolicy(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) DeleteRetryPolicy(id string) error {
	return s.store.DeleteRetryPolicy(id)
}
