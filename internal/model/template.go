package model

import (
	"strings"
)

const (
	TemplateActive   = "active"
	TemplateDisabled = "disabled"
)

type Template struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	FlowID      string `json:"flow_id"`
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func (t *Template) Validate() error {
	t.Name = strings.TrimSpace(t.Name)
	t.Content = strings.TrimSpace(t.Content)
	if t.Name == "" {
		return NewValidationError("name", "模板名称不能为空")
	}
	if t.Content == "" {
		return NewValidationError("content", "模板内容不能为空")
	}
	if t.ContentType == "" {
		t.ContentType = "application/json"
	}
	if t.Status == "" {
		t.Status = TemplateActive
	}
	if t.Status != TemplateActive && t.Status != TemplateDisabled {
		return NewValidationError("status", "状态不合法")
	}
	return nil
}

func (t *Template) ReplaceVars(vars map[string]string) string {
	result := t.Content
	for k, v := range vars {
		result = strings.ReplaceAll(result, "{{"+k+"}}", v)
	}
	return result
}

type TemplateFilter struct {
	FlowID string
	Status string
	Name   string
}

func (f TemplateFilter) Match(t *Template) bool {
	if f.FlowID != "" && t.FlowID != f.FlowID {
		return false
	}
	if f.Status != "" && t.Status != f.Status {
		return false
	}
	if f.Name != "" && t.Name != f.Name {
		return false
	}
	return true
}
