package model

import (
	"strings"
)

type Dependency struct {
	ID              string `json:"id"`
	FlowID          string `json:"flow_id"`
	StepID          string `json:"step_id"`
	DependsOnStepID string `json:"depends_on_step_id"`
	Condition       string `json:"condition"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

func (d *Dependency) Validate() error {
	if d.FlowID == "" {
		return NewValidationError("flow_id", "流程ID不能为空")
	}
	if d.StepID == "" {
		return NewValidationError("step_id", "步骤ID不能为空")
	}
	if d.DependsOnStepID == "" {
		return NewValidationError("depends_on_step_id", "依赖步骤ID不能为空")
	}
	if d.StepID == d.DependsOnStepID {
		return NewValidationError("depends_on_step_id", "步骤不能依赖自身")
	}
	d.Condition = strings.TrimSpace(d.Condition)
	return nil
}

type DependencyFilter struct {
	FlowID          string
	StepID          string
	DependsOnStepID string
}

func (f DependencyFilter) Match(d *Dependency) bool {
	if f.FlowID != "" && d.FlowID != f.FlowID {
		return false
	}
	if f.StepID != "" && d.StepID != f.StepID {
		return false
	}
	if f.DependsOnStepID != "" && d.DependsOnStepID != f.DependsOnStepID {
		return false
	}
	return true
}
