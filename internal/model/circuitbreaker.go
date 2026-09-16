package model

const (
	CircuitBreakerActive   = "active"
	CircuitBreakerDisabled = "disabled"
)

type CircuitBreaker struct {
	ID         string `json:"id"`
	ServiceID  string `json:"service_id"`
	Threshold  int    `json:"threshold"`
	CooldownMs int    `json:"cooldown_ms"`
	Status     string `json:"status"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

func (c *CircuitBreaker) Validate() error {
	if c.ServiceID == "" {
		return NewValidationError("service_id", "服务ID不能为空")
	}
	if c.Threshold <= 0 {
		return NewValidationError("threshold", "阈值必须大于0")
	}
	if c.CooldownMs <= 0 {
		return NewValidationError("cooldown_ms", "冷却时间必须大于0")
	}
	if c.Status == "" {
		c.Status = CircuitBreakerActive
	}
	if c.Status != CircuitBreakerActive && c.Status != CircuitBreakerDisabled {
		return NewValidationError("status", "状态不合法")
	}
	return nil
}

type CircuitBreakerFilter struct {
	ServiceID string
	Status    string
}

func (f CircuitBreakerFilter) Match(c *CircuitBreaker) bool {
	if f.ServiceID != "" && c.ServiceID != f.ServiceID {
		return false
	}
	if f.Status != "" && c.Status != f.Status {
		return false
	}
	return true
}
