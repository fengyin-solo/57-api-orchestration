package store

import (
	"orchestr/internal/model"
)

func (s *MemoryStore) CreateRetryPolicy(r *model.RetryPolicy) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.retryPolicies {
		if exist.Name == r.Name {
			return ErrConflict
		}
	}
	s.retryPolicies[r.ID] = r
	return nil
}

func (s *MemoryStore) GetRetryPolicy(id string) (*model.RetryPolicy, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.retryPolicies[id]
	if !ok {
		return nil, ErrNotFound
	}
	return r, nil
}

func (s *MemoryStore) GetRetryPolicyByName(name string) (*model.RetryPolicy, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, r := range s.retryPolicies {
		if r.Name == name {
			return r, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListRetryPolicies() []*model.RetryPolicy {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.RetryPolicy, 0, len(s.retryPolicies))
	for _, r := range s.retryPolicies {
		list = append(list, r)
	}
	return list
}

func (s *MemoryStore) UpdateRetryPolicy(r *model.RetryPolicy) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.retryPolicies[r.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.retryPolicies {
		if exist.ID != r.ID && exist.Name == r.Name {
			return ErrConflict
		}
	}
	s.retryPolicies[r.ID] = r
	return nil
}

func (s *MemoryStore) DeleteRetryPolicy(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.retryPolicies[id]; !ok {
		return ErrNotFound
	}
	delete(s.retryPolicies, id)
	return nil
}
