package store

import (
	"orchestr/internal/model"
)

func (s *MemoryStore) CreateWebhook(w *model.Webhook) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.webhooks[w.ID] = w
	return nil
}

func (s *MemoryStore) GetWebhook(id string) (*model.Webhook, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	w, ok := s.webhooks[id]
	if !ok {
		return nil, ErrNotFound
	}
	return w, nil
}

func (s *MemoryStore) ListWebhooks() []*model.Webhook {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Webhook, 0, len(s.webhooks))
	for _, w := range s.webhooks {
		list = append(list, w)
	}
	return list
}

func (s *MemoryStore) UpdateWebhook(w *model.Webhook) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.webhooks[w.ID]; !ok {
		return ErrNotFound
	}
	s.webhooks[w.ID] = w
	return nil
}

func (s *MemoryStore) DeleteWebhook(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.webhooks[id]; !ok {
		return ErrNotFound
	}
	delete(s.webhooks, id)
	return nil
}

func (s *MemoryStore) ListWebhooksByFlowID(flowID string) []*model.Webhook {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Webhook, 0)
	for _, w := range s.webhooks {
		if w.FlowID == flowID {
			list = append(list, w)
		}
	}
	return list
}
