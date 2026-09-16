package model

import (
	"strings"
)

const (
	WebhookActive   = "active"
	WebhookDisabled = "disabled"
)

type Webhook struct {
	ID          string `json:"id"`
	FlowID      string `json:"flow_id"`
	URL         string `json:"url"`
	Method      string `json:"method"`
	Headers     string `json:"headers"`
	RetryCount  int    `json:"retry_count"`
	TimeoutMs   int    `json:"timeout_ms"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func (w *Webhook) Validate() error {
	w.URL = strings.TrimSpace(w.URL)
	w.Method = strings.TrimSpace(w.Method)
	if w.FlowID == "" {
		return NewValidationError("flow_id", "流程ID不能为空")
	}
	if w.URL == "" {
		return NewValidationError("url", "Webhook URL 不能为空")
	}
	if w.Method == "" {
		w.Method = "POST"
	}
	if w.TimeoutMs <= 0 {
		w.TimeoutMs = 5000
	}
	if w.RetryCount < 0 {
		w.RetryCount = 0
	}
	if w.Status == "" {
		w.Status = WebhookActive
	}
	if w.Status != WebhookActive && w.Status != WebhookDisabled {
		return NewValidationError("status", "状态不合法")
	}
	return nil
}

type WebhookFilter struct {
	FlowID string
	Status string
}

func (f WebhookFilter) Match(w *Webhook) bool {
	if f.FlowID != "" && w.FlowID != f.FlowID {
		return false
	}
	if f.Status != "" && w.Status != f.Status {
		return false
	}
	return true
}
