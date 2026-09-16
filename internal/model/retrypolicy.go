package model

import (
	"strings"
)

const (
	BackoffFixed     = "fixed"
	BackoffLinear    = "linear"
	BackoffExponential = "exponential"
	RetryPolicyActive   = "active"
	RetryPolicyDisabled = "disabled"
)

type RetryPolicy struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	MaxRetries     int    `json:"max_retries"`
	BackoffType    string `json:"backoff_type"`
	InitialDelayMs int    `json:"initial_delay_ms"`
	Status         string `json:"status"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

func (r *RetryPolicy) Validate() error {
	r.Name = strings.TrimSpace(r.Name)
	if r.Name == "" {
		return NewValidationError("name", "策略名称不能为空")
	}
	if r.MaxRetries < 0 {
		return NewValidationError("max_retries", "最大重试次数不能为负")
	}
	if r.BackoffType == "" {
		r.BackoffType = BackoffFixed
	}
	if r.BackoffType != BackoffFixed && r.BackoffType != BackoffLinear && r.BackoffType != BackoffExponential {
		return NewValidationError("backoff_type", "退避类型不合法")
	}
	if r.InitialDelayMs <= 0 {
		r.InitialDelayMs = 1000
	}
	if r.Status == "" {
		r.Status = RetryPolicyActive
	}
	if r.Status != RetryPolicyActive && r.Status != RetryPolicyDisabled {
		return NewValidationError("status", "状态不合法")
	}
	return nil
}

func (r *RetryPolicy) CalculateDelay(attempt int) int {
	if attempt < 0 {
		attempt = 0
	}
	switch r.BackoffType {
	case BackoffLinear:
		return r.InitialDelayMs * (attempt + 1)
	case BackoffExponential:
		delay := r.InitialDelayMs
		for i := 0; i < attempt; i++ {
			delay *= 2
		}
		return delay
	default:
		return r.InitialDelayMs
	}
}

type RetryPolicyFilter struct {
	Status  string
	Keyword string
}

func (f RetryPolicyFilter) Match(r *RetryPolicy) bool {
	if f.Status != "" && r.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(r.Name), k) {
			return false
		}
	}
	return true
}
