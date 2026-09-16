package model

import (
	"encoding/json"
	"strings"
)

const (
	StepActive   = "active"
	StepDisabled = "disabled"
)

type Step struct {
	ID              string          `json:"id"`
	FlowID          string          `json:"flow_id"`
	ServiceID       string          `json:"service_id"`
	Method          string          `json:"method"`
	Path            string          `json:"path"`
	ParamsMapping   json.RawMessage `json:"params_mapping"`
	Order           int             `json:"order"`
	TimeoutMs       int             `json:"timeout_ms"`
	RetryCount      int             `json:"retry_count"`
	Status          string          `json:"status"`
	CreatedAt       string          `json:"created_at"`
	UpdatedAt       string          `json:"updated_at"`
}

func (s *Step) Validate() error {
	s.Method = strings.TrimSpace(s.Method)
	s.Path = strings.TrimSpace(s.Path)
	if s.FlowID == "" {
		return NewValidationError("flow_id", "流程ID不能为空")
	}
	if s.ServiceID == "" {
		return NewValidationError("service_id", "服务ID不能为空")
	}
	if s.Method == "" {
		return NewValidationError("method", "HTTP方法不能为空")
	}
	if s.Path == "" {
		return NewValidationError("path", "请求路径不能为空")
	}
	if s.TimeoutMs <= 0 {
		s.TimeoutMs = 5000
	}
	if s.RetryCount < 0 {
		s.RetryCount = 0
	}
	if s.Status == "" {
		s.Status = StepActive
	}
	if s.Status != StepActive && s.Status != StepDisabled {
		return NewValidationError("status", "步骤状态不合法")
	}
	return nil
}

type StepFilter struct {
	FlowID   string
	Status   string
	Method   string
	ServiceID string
}

func (f StepFilter) Match(s *Step) bool {
	if f.FlowID != "" && s.FlowID != f.FlowID {
		return false
	}
	if f.Status != "" && s.Status != f.Status {
		return false
	}
	if f.Method != "" && s.Method != f.Method {
		return false
	}
	if f.ServiceID != "" && s.ServiceID != f.ServiceID {
		return false
	}
	return true
}
