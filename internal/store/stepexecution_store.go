package store

import (
	"orchestr/internal/model"
)

func (s *MemoryStore) CreateStepExecution(se *model.StepExecution) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stepExecutions[se.ID] = se
	return nil
}

func (s *MemoryStore) GetStepExecution(id string) (*model.StepExecution, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	se, ok := s.stepExecutions[id]
	if !ok {
		return nil, ErrNotFound
	}
	return se, nil
}

func (s *MemoryStore) ListStepExecutions() []*model.StepExecution {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.StepExecution, 0, len(s.stepExecutions))
	for _, se := range s.stepExecutions {
		list = append(list, se)
	}
	return list
}

func (s *MemoryStore) UpdateStepExecution(se *model.StepExecution) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.stepExecutions[se.ID]; !ok {
		return ErrNotFound
	}
	s.stepExecutions[se.ID] = se
	return nil
}

func (s *MemoryStore) DeleteStepExecution(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.stepExecutions[id]; !ok {
		return ErrNotFound
	}
	delete(s.stepExecutions, id)
	return nil
}

func (s *MemoryStore) ListStepExecutionsByExecutionID(executionID string) []*model.StepExecution {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.StepExecution, 0)
	for _, se := range s.stepExecutions {
		if se.ExecutionID == executionID {
			list = append(list, se)
		}
	}
	return list
}
