package model

import (
	"encoding/json"
)

const (
	StepExecutionPending   = "pending"
	StepExecutionRunning   = "running"
	StepExecutionCompleted = "completed"
	StepExecutionFailed    = "failed"
	StepExecutionSkipped   = "skipped"
)

var stepExecutionTransitions = map[string]map[string]bool{
	StepExecutionPending:   {StepExecutionRunning: true, StepExecutionSkipped: true},
	StepExecutionRunning:   {StepExecutionCompleted: true, StepExecutionFailed: true, StepExecutionSkipped: true},
	StepExecutionCompleted: {},
	StepExecutionFailed:    {},
	StepExecutionSkipped:   {},
}

func StepExecutionCanTransition(from, to string) bool {
	if m, ok := stepExecutionTransitions[from]; ok {
		return m[to]
	}
	return false
}

type StepExecution struct {
	ID          string          `json:"id"`
	ExecutionID string          `json:"execution_id"`
	StepID      string          `json:"step_id"`
	Status      string          `json:"status"`
	Response    json.RawMessage `json:"response"`
	DurationMs  int             `json:"duration_ms"`
	StartedAt   string          `json:"started_at"`
	FinishedAt  string          `json:"finished_at"`
}

func (s *StepExecution) Validate() error {
	if s.ExecutionID == "" {
		return NewValidationError("execution_id", "执行实例ID不能为空")
	}
	if s.StepID == "" {
		return NewValidationError("step_id", "步骤ID不能为空")
	}
	if s.Status == "" {
		s.Status = StepExecutionPending
	}
	if s.Status != StepExecutionPending && s.Status != StepExecutionRunning && s.Status != StepExecutionCompleted && s.Status != StepExecutionFailed && s.Status != StepExecutionSkipped {
		return NewValidationError("status", "步骤执行状态不合法")
	}
	return nil
}

type StepExecutionFilter struct {
	ExecutionID string
	Status      string
}

func (f StepExecutionFilter) Match(s *StepExecution) bool {
	if f.ExecutionID != "" && s.ExecutionID != f.ExecutionID {
		return false
	}
	if f.Status != "" && s.Status != f.Status {
		return false
	}
	return true
}
