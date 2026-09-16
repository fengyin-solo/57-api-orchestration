package store

import (
	"orchestr/internal/model"
)

func (s *MemoryStore) CreateEnvVar(e *model.EnvVar) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.envVars[e.ID] = e
	return nil
}

func (s *MemoryStore) GetEnvVar(id string) (*model.EnvVar, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.envVars[id]
	if !ok {
		return nil, ErrNotFound
	}
	return e, nil
}

func (s *MemoryStore) ListEnvVars() []*model.EnvVar {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.EnvVar, 0, len(s.envVars))
	for _, e := range s.envVars {
		list = append(list, e)
	}
	return list
}

func (s *MemoryStore) UpdateEnvVar(e *model.EnvVar) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.envVars[e.ID]; !ok {
		return ErrNotFound
	}
	s.envVars[e.ID] = e
	return nil
}

func (s *MemoryStore) DeleteEnvVar(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.envVars[id]; !ok {
		return ErrNotFound
	}
	delete(s.envVars, id)
	return nil
}

func (s *MemoryStore) ListEnvVarsByFlowID(flowID string) []*model.EnvVar {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.EnvVar, 0)
	for _, e := range s.envVars {
		if e.FlowID == flowID {
			list = append(list, e)
		}
	}
	return list
}
