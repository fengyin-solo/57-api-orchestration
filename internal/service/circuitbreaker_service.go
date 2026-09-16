package service

import (
	"sort"
	"time"

	"orchestr/internal/model"
	"orchestr/pkg/idgen"
)

func (s *Service) CreateCircuitBreaker(input model.CircuitBreaker) (*model.CircuitBreaker, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetService(input.ServiceID); err != nil {
		return nil, model.NewValidationError("service_id", "服务不存在")
	}
	now := time.Now().Format(time.RFC3339)
	c := &model.CircuitBreaker{
		ID:         idgen.Hex(),
		ServiceID:  input.ServiceID,
		Threshold:  input.Threshold,
		CooldownMs: input.CooldownMs,
		Status:     input.Status,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := s.store.CreateCircuitBreaker(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) GetCircuitBreaker(id string) (*model.CircuitBreaker, error) {
	return s.store.GetCircuitBreaker(id)
}

func (s *Service) ListCircuitBreakers(filter model.CircuitBreakerFilter, page, size int) ([]*model.CircuitBreaker, int, error) {
	all := s.store.ListCircuitBreakers()
	matched := make([]*model.CircuitBreaker, 0, len(all))
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
		return []*model.CircuitBreaker{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateCircuitBreaker(id string, input model.CircuitBreaker) (*model.CircuitBreaker, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	c, err := s.store.GetCircuitBreaker(id)
	if err != nil {
		return nil, err
	}
	c.ServiceID = input.ServiceID
	c.Threshold = input.Threshold
	c.CooldownMs = input.CooldownMs
	c.Status = input.Status
	c.UpdatedAt = time.Now().Format(time.RFC3339)
	if err := s.store.UpdateCircuitBreaker(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) DeleteCircuitBreaker(id string) error {
	return s.store.DeleteCircuitBreaker(id)
}
