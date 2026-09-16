package store

import (
	"orchestr/internal/model"
)

func (s *MemoryStore) CreateExecution(e *model.Execution) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.executions[e.ID] = e
	return nil
}

func (s *MemoryStore) GetExecution(id string) (*model.Execution, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.executions[id]
	if !ok {
		return nil, ErrNotFound
	}
	return e, nil
}

func (s *MemoryStore) ListExecutions() []*model.Execution {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Execution, 0, len(s.executions))
	for _, e := range s.executions {
		list = append(list, e)
	}
	return list
}

func (s *MemoryStore) UpdateExecution(e *model.Execution) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.executions[e.ID]; !ok {
		return ErrNotFound
	}
	s.executions[e.ID] = e
	return nil
}

func (s *MemoryStore) DeleteExecution(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.executions[id]; !ok {
		return ErrNotFound
	}
	delete(s.executions, id)
	return nil
}
