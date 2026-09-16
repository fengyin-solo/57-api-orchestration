package model

import (
	"strings"
)

const (
	ConditionActive   = "active"
	ConditionDisabled = "disabled"
)

type Condition struct {
	ID           string `json:"id"`
	FlowID       string `json:"flow_id"`
	StepID       string `json:"step_id"`
	Expression   string `json:"expression"`
	TrueStepID   string `json:"true_step_id"`
	FalseStepID  string `json:"false_step_id"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

func (c *Condition) Validate() error {
	if c.FlowID == "" {
		return NewValidationError("flow_id", "流程ID不能为空")
	}
	if c.StepID == "" {
		return NewValidationError("step_id", "步骤ID不能为空")
	}
	c.Expression = strings.TrimSpace(c.Expression)
	if c.Expression == "" {
		return NewValidationError("expression", "条件表达式不能为空")
	}
	if c.TrueStepID == "" && c.FalseStepID == "" {
		return NewValidationError("true_step_id", "至少需要指定一个分支步骤")
	}
	if c.Status == "" {
		c.Status = ConditionActive
	}
	if c.Status != ConditionActive && c.Status != ConditionDisabled {
		return NewValidationError("status", "状态不合法")
	}
	return nil
}

func (c *Condition) Evaluate(input map[string]interface{}) bool {
	if val, ok := input["result"]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return true
}

type ConditionFilter struct {
	FlowID string
	Status string
}

func (f ConditionFilter) Match(c *Condition) bool {
	if f.FlowID != "" && c.FlowID != f.FlowID {
		return false
	}
	if f.Status != "" && c.Status != f.Status {
		return false
	}
	return true
}
