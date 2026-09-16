package store

import (
	"orchestr/internal/model"
)

func (s *MemoryStore) CreateCircuitBreaker(c *model.CircuitBreaker) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.circuitBreakers {
		if exist.ServiceID == c.ServiceID {
			return ErrConflict
		}
	}
	s.circuitBreakers[c.ID] = c
	return nil
}

func (s *MemoryStore) GetCircuitBreaker(id string) (*model.CircuitBreaker, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.circuitBreakers[id]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
}

func (s *MemoryStore) GetCircuitBreakerByServiceID(serviceID string) (*model.CircuitBreaker, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, c := range s.circuitBreakers {
		if c.ServiceID == serviceID {
			return c, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListCircuitBreakers() []*model.CircuitBreaker {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.CircuitBreaker, 0, len(s.circuitBreakers))
	for _, c := range s.circuitBreakers {
		list = append(list, c)
	}
	return list
}

func (s *MemoryStore) UpdateCircuitBreaker(c *model.CircuitBreaker) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.circuitBreakers[c.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.circuitBreakers {
		if exist.ID != c.ID && exist.ServiceID == c.ServiceID {
			return ErrConflict
		}
	}
	s.circuitBreakers[c.ID] = c
	return nil
}

func (s *MemoryStore) DeleteCircuitBreaker(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.circuitBreakers[id]; !ok {
		return ErrNotFound
	}
	delete(s.circuitBreakers, id)
	return nil
}
