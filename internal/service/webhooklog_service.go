package service

import (
	"sort"
	"time"

	"orchestr/internal/model"
	"orchestr/pkg/idgen"
)

func (s *Service) CreateWebhookLog(input model.WebhookLog) (*model.WebhookLog, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetWebhook(input.WebhookID); err != nil {
		return nil, model.NewValidationError("webhook_id", "Webhook 不存在")
	}
	input.ID = idgen.Hex()
	input.CreatedAt = time.Now().Format(time.RFC3339)
	if err := s.store.CreateWebhookLog(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *Service) GetWebhookLog(id string) (*model.WebhookLog, error) {
	return s.store.GetWebhookLog(id)
}

func (s *Service) ListWebhookLogs(filter model.WebhookLogFilter, page, size int) ([]*model.WebhookLog, int, error) {
	all := s.store.ListWebhookLogs()
	matched := make([]*model.WebhookLog, 0, len(all))
	for _, w := range all {
		if filter.Match(w) {
			matched = append(matched, w)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt > matched[j].CreatedAt
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.WebhookLog{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateWebhookLog(id string, input model.WebhookLog) (*model.WebhookLog, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	w, err := s.store.GetWebhookLog(id)
	if err != nil {
		return nil, err
	}
	w.WebhookID = input.WebhookID
	w.Payload = input.Payload
	w.Status = input.Status
	w.Response = input.Response
	if err := s.store.UpdateWebhookLog(w); err != nil {
		return nil, err
	}
	return w, nil
}

func (s *Service) DeleteWebhookLog(id string) error {
	return s.store.DeleteWebhookLog(id)
}
