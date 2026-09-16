package service

import (
	"sort"
	"time"

	"orchestr/internal/model"
	"orchestr/pkg/idgen"
)

func (s *Service) CreateWebhook(input model.Webhook) (*model.Webhook, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetFlow(input.FlowID); err != nil {
		return nil, model.NewValidationError("flow_id", "流程不存在")
	}
	now := time.Now().Format(time.RFC3339)
	w := &model.Webhook{
		ID:         idgen.Hex(),
		FlowID:     input.FlowID,
		URL:        input.URL,
		Method:     input.Method,
		Headers:    input.Headers,
		RetryCount: input.RetryCount,
		TimeoutMs:  input.TimeoutMs,
		Status:     input.Status,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := s.store.CreateWebhook(w); err != nil {
		return nil, err
	}
	return w, nil
}

func (s *Service) GetWebhook(id string) (*model.Webhook, error) {
	return s.store.GetWebhook(id)
}

func (s *Service) ListWebhooks(filter model.WebhookFilter, page, size int) ([]*model.Webhook, int, error) {
	all := s.store.ListWebhooks()
	matched := make([]*model.Webhook, 0, len(all))
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
		return []*model.Webhook{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateWebhook(id string, input model.Webhook) (*model.Webhook, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	w, err := s.store.GetWebhook(id)
	if err != nil {
		return nil, err
	}
	w.FlowID = input.FlowID
	w.URL = input.URL
	w.Method = input.Method
	w.Headers = input.Headers
	w.RetryCount = input.RetryCount
	w.TimeoutMs = input.TimeoutMs
	w.Status = input.Status
	w.UpdatedAt = time.Now().Format(time.RFC3339)
	if err := s.store.UpdateWebhook(w); err != nil {
		return nil, err
	}
	return w, nil
}

func (s *Service) DeleteWebhook(id string) error {
	return s.store.DeleteWebhook(id)
}

func (s *Service) TriggerWebhooks(flowID string, payload string) error {
	webhooks := s.store.ListWebhooksByFlowID(flowID)
	for _, w := range webhooks {
		if w.Status != model.WebhookActive {
			continue
		}
		s.log.Infof("触发 webhook %s %s", w.Method, w.URL)
	}
	return nil
}
