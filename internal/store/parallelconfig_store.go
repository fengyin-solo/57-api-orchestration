package store

import (
	"orchestr/internal/model"
)

func (s *MemoryStore) CreateParallelConfig(p *model.ParallelConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.parallelConfigs {
		if exist.FlowID == p.FlowID {
			return ErrConflict
		}
	}
	s.parallelConfigs[p.ID] = p
	return nil
}

func (s *MemoryStore) GetParallelConfig(id string) (*model.ParallelConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.parallelConfigs[id]
	if !ok {
		return nil, ErrNotFound
	}
	return p, nil
}

func (s *MemoryStore) ListParallelConfigs() []*model.ParallelConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.ParallelConfig, 0, len(s.parallelConfigs))
	for _, p := range s.parallelConfigs {
		list = append(list, p)
	}
	return list
}

func (s *MemoryStore) UpdateParallelConfig(p *model.ParallelConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.parallelConfigs[p.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.parallelConfigs {
		if exist.ID != p.ID && exist.FlowID == p.FlowID {
			return ErrConflict
		}
	}
	s.parallelConfigs[p.ID] = p
	return nil
}

func (s *MemoryStore) DeleteParallelConfig(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.parallelConfigs[id]; !ok {
		return ErrNotFound
	}
	delete(s.parallelConfigs, id)
	return nil
}

func (s *MemoryStore) GetParallelConfigByFlowID(flowID string) (*model.ParallelConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, p := range s.parallelConfigs {
		if p.FlowID == flowID {
			return p, nil
		}
	}
	return nil, ErrNotFound
}
