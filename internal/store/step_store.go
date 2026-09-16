package store

import (
	"orchestr/internal/model"
)

func (s *MemoryStore) CreateStep(st *model.Step) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.steps[st.ID] = st
	return nil
}

func (s *MemoryStore) GetStep(id string) (*model.Step, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	st, ok := s.steps[id]
	if !ok {
		return nil, ErrNotFound
	}
	return st, nil
}

func (s *MemoryStore) ListSteps() []*model.Step {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Step, 0, len(s.steps))
	for _, st := range s.steps {
		list = append(list, st)
	}
	return list
}

func (s *MemoryStore) UpdateStep(st *model.Step) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.steps[st.ID]; !ok {
		return ErrNotFound
	}
	s.steps[st.ID] = st
	return nil
}

func (s *MemoryStore) DeleteStep(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.steps[id]; !ok {
		return ErrNotFound
	}
	delete(s.steps, id)
	return nil
}

func (s *MemoryStore) ListStepsByFlowID(flowID string) []*model.Step {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Step, 0)
	for _, st := range s.steps {
		if st.FlowID == flowID {
			list = append(list, st)
		}
	}
	return list
}
