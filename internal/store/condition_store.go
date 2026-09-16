package store

import (
	"orchestr/internal/model"
)

func (s *MemoryStore) CreateCondition(c *model.Condition) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.conditions[c.ID] = c
	return nil
}

func (s *MemoryStore) GetCondition(id string) (*model.Condition, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.conditions[id]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
}

func (s *MemoryStore) ListConditions() []*model.Condition {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Condition, 0, len(s.conditions))
	for _, c := range s.conditions {
		list = append(list, c)
	}
	return list
}

func (s *MemoryStore) UpdateCondition(c *model.Condition) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.conditions[c.ID]; !ok {
		return ErrNotFound
	}
	s.conditions[c.ID] = c
	return nil
}

func (s *MemoryStore) DeleteCondition(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.conditions[id]; !ok {
		return ErrNotFound
	}
	delete(s.conditions, id)
	return nil
}

func (s *MemoryStore) ListConditionsByFlowID(flowID string) []*model.Condition {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Condition, 0)
	for _, c := range s.conditions {
		if c.FlowID == flowID {
			list = append(list, c)
		}
	}
	return list
}
