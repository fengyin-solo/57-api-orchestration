package store

import (
	"orchestr/internal/model"
)

func (s *MemoryStore) CreateWebhookLog(w *model.WebhookLog) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.webhookLogs[w.ID] = w
	return nil
}

func (s *MemoryStore) GetWebhookLog(id string) (*model.WebhookLog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	w, ok := s.webhookLogs[id]
	if !ok {
		return nil, ErrNotFound
	}
	return w, nil
}

func (s *MemoryStore) ListWebhookLogs() []*model.WebhookLog {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.WebhookLog, 0, len(s.webhookLogs))
	for _, w := range s.webhookLogs {
		list = append(list, w)
	}
	return list
}

func (s *MemoryStore) UpdateWebhookLog(w *model.WebhookLog) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.webhookLogs[w.ID]; !ok {
		return ErrNotFound
	}
	s.webhookLogs[w.ID] = w
	return nil
}

func (s *MemoryStore) DeleteWebhookLog(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.webhookLogs[id]; !ok {
		return ErrNotFound
	}
	delete(s.webhookLogs, id)
	return nil
}

func (s *MemoryStore) ListWebhookLogsByWebhookID(webhookID string) []*model.WebhookLog {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.WebhookLog, 0)
	for _, w := range s.webhookLogs {
		if w.WebhookID == webhookID {
			list = append(list, w)
		}
	}
	return list
}
