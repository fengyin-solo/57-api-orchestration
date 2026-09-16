package model

import (
	"strings"
)

const (
	WebhookLogPending   = "pending"
	WebhookLogDelivered = "delivered"
	WebhookLogFailed    = "failed"
)

type WebhookLog struct {
	ID        string `json:"id"`
	WebhookID string `json:"webhook_id"`
	Payload   string `json:"payload"`
	Status    string `json:"status"`
	Response  string `json:"response"`
	CreatedAt string `json:"created_at"`
}

func (w *WebhookLog) Validate() error {
	if w.WebhookID == "" {
		return NewValidationError("webhook_id", "Webhook ID 不能为空")
	}
	w.Payload = strings.TrimSpace(w.Payload)
	if w.Status == "" {
		w.Status = WebhookLogPending
	}
	if w.Status != WebhookLogPending && w.Status != WebhookLogDelivered && w.Status != WebhookLogFailed {
		return NewValidationError("status", "状态不合法")
	}
	return nil
}

type WebhookLogFilter struct {
	WebhookID string
	Status    string
}

func (f WebhookLogFilter) Match(w *WebhookLog) bool {
	if f.WebhookID != "" && w.WebhookID != f.WebhookID {
		return false
	}
	if f.Status != "" && w.Status != f.Status {
		return false
	}
	return true
}
