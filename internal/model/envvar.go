package model

import (
	"strings"
)

const (
	EnvVarActive   = "active"
	EnvVarDisabled = "disabled"
)

type EnvVar struct {
	ID        string `json:"id"`
	FlowID    string `json:"flow_id"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	IsSecret  bool   `json:"is_secret"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (e *EnvVar) Validate() error {
	e.Key = strings.TrimSpace(e.Key)
	if e.Key == "" {
		return NewValidationError("key", "变量键不能为空")
	}
	if e.FlowID == "" {
		return NewValidationError("flow_id", "流程ID不能为空")
	}
	if e.Status == "" {
		e.Status = EnvVarActive
	}
	if e.Status != EnvVarActive && e.Status != EnvVarDisabled {
		return NewValidationError("status", "状态不合法")
	}
	return nil
}

type EnvVarFilter struct {
	FlowID string
	Status string
	Key    string
}

func (f EnvVarFilter) Match(e *EnvVar) bool {
	if f.FlowID != "" && e.FlowID != f.FlowID {
		return false
	}
	if f.Status != "" && e.Status != f.Status {
		return false
	}
	if f.Key != "" && e.Key != f.Key {
		return false
	}
	return true
}
