package model

import (
	"strings"
)

const (
	ServiceAuthNone  = "none"
	ServiceAuthBasic = "basic"
	ServiceAuthToken = "token"
	ServiceActive    = "active"
	ServiceDisabled  = "disabled"
)

type Service struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	BaseURL     string `json:"base_url"`
	AuthType    string `json:"auth_type"`
	TimeoutMs   int    `json:"timeout_ms"`
	RetryCount  int    `json:"retry_count"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func (s *Service) Validate() error {
	s.Name = strings.TrimSpace(s.Name)
	s.BaseURL = strings.TrimSpace(s.BaseURL)
	if s.Name == "" {
		return NewValidationError("name", "服务名称不能为空")
	}
	if s.BaseURL == "" {
		return NewValidationError("base_url", "BaseURL 不能为空")
	}
	if s.AuthType == "" {
		s.AuthType = ServiceAuthNone
	}
	if s.AuthType != ServiceAuthNone && s.AuthType != ServiceAuthBasic && s.AuthType != ServiceAuthToken {
		return NewValidationError("auth_type", "认证类型不合法")
	}
	if s.TimeoutMs <= 0 {
		s.TimeoutMs = 5000
	}
	if s.RetryCount < 0 {
		s.RetryCount = 0
	}
	if s.Status == "" {
		s.Status = ServiceActive
	}
	if s.Status != ServiceActive && s.Status != ServiceDisabled {
		return NewValidationError("status", "状态不合法")
	}
	return nil
}

type ServiceFilter struct {
	Status  string
	Keyword string
}

func (f ServiceFilter) Match(s *Service) bool {
	if f.Status != "" && s.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(s.Name), k) {
			return false
		}
	}
	return true
}
