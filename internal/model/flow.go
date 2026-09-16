package model

import (
	"strings"
)

const (
	FlowActive   = "active"
	FlowDisabled = "disabled"
)

type Flow struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func (f *Flow) Validate() error {
	f.Name = strings.TrimSpace(f.Name)
	if f.Name == "" {
		return NewValidationError("name", "流程名称不能为空")
	}
	if f.Status == "" {
		f.Status = FlowActive
	}
	if f.Status != FlowActive && f.Status != FlowDisabled {
		return NewValidationError("status", "流程状态不合法")
	}
	return nil
}

type FlowFilter struct {
	Status  string
	Keyword string
}

func (ff FlowFilter) Match(f *Flow) bool {
	if ff.Status != "" && f.Status != ff.Status {
		return false
	}
	if ff.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(ff.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(f.Name), k) {
			return false
		}
	}
	return true
}
