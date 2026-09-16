package store

import (
	"orchestr/internal/model"
)

func (s *MemoryStore) CreateDependency(d *model.Dependency) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.dependencies[d.ID] = d
	return nil
}

func (s *MemoryStore) GetDependency(id string) (*model.Dependency, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	d, ok := s.dependencies[id]
	if !ok {
		return nil, ErrNotFound
	}
	return d, nil
}

func (s *MemoryStore) ListDependencies() []*model.Dependency {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Dependency, 0, len(s.dependencies))
	for _, d := range s.dependencies {
		list = append(list, d)
	}
	return list
}

func (s *MemoryStore) UpdateDependency(d *model.Dependency) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.dependencies[d.ID]; !ok {
		return ErrNotFound
	}
	s.dependencies[d.ID] = d
	return nil
}

func (s *MemoryStore) DeleteDependency(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.dependencies[id]; !ok {
		return ErrNotFound
	}
	delete(s.dependencies, id)
	return nil
}

func (s *MemoryStore) ListDependenciesByFlowID(flowID string) []*model.Dependency {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Dependency, 0)
	for _, d := range s.dependencies {
		if d.FlowID == flowID {
			list = append(list, d)
		}
	}
	return list
}
