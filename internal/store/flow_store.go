package store

import (
	"orchestr/internal/model"
)

func (s *MemoryStore) CreateFlow(f *model.Flow) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.flows {
		if exist.Name == f.Name {
			return ErrConflict
		}
	}
	s.flows[f.ID] = f
	return nil
}

func (s *MemoryStore) GetFlow(id string) (*model.Flow, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	f, ok := s.flows[id]
	if !ok {
		return nil, ErrNotFound
	}
	return f, nil
}

func (s *MemoryStore) GetFlowByName(name string) (*model.Flow, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, f := range s.flows {
		if f.Name == name {
			return f, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListFlows() []*model.Flow {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Flow, 0, len(s.flows))
	for _, f := range s.flows {
		list = append(list, f)
	}
	return list
}

func (s *MemoryStore) UpdateFlow(f *model.Flow) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.flows[f.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.flows {
		if exist.ID != f.ID && exist.Name == f.Name {
			return ErrConflict
		}
	}
	s.flows[f.ID] = f
	return nil
}

func (s *MemoryStore) DeleteFlow(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.flows[id]; !ok {
		return ErrNotFound
	}
	delete(s.flows, id)
	return nil
}
