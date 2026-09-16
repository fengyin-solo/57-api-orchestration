package model

import (
	"encoding/json"
)

const (
	ExecutionPending   = "pending"
	ExecutionRunning   = "running"
	ExecutionCompleted = "completed"
	ExecutionFailed    = "failed"
	ExecutionTimeout   = "timeout"
)

var executionTransitions = map[string]map[string]bool{
	ExecutionPending:   {ExecutionRunning: true},
	ExecutionRunning:   {ExecutionCompleted: true, ExecutionFailed: true, ExecutionTimeout: true},
	ExecutionCompleted: {},
	ExecutionFailed:    {},
	ExecutionTimeout:   {},
}

func ExecutionCanTransition(from, to string) bool {
	if m, ok := executionTransitions[from]; ok {
		return m[to]
	}
	return false
}

type Execution struct {
	ID         string          `json:"id"`
	FlowID     string          `json:"flow_id"`
	Status     string          `json:"status"`
	Input      json.RawMessage `json:"input"`
	Output     json.RawMessage `json:"output"`
	ErrorMsg   string          `json:"error_msg"`
	CreatedAt  string          `json:"created_at"`
	StartedAt  string          `json:"started_at"`
	FinishedAt string          `json:"finished_at"`
}

func (e *Execution) Validate() error {
	if e.FlowID == "" {
		return NewValidationError("flow_id", "流程ID不能为空")
	}
	if e.Status == "" {
		e.Status = ExecutionPending
	}
	if e.Status != ExecutionPending && e.Status != ExecutionRunning && e.Status != ExecutionCompleted && e.Status != ExecutionFailed && e.Status != ExecutionTimeout {
		return NewValidationError("status", "执行状态不合法")
	}
	return nil
}

type ExecutionFilter struct {
	FlowID string
	Status string
}

func (f ExecutionFilter) Match(e *Execution) bool {
	if f.FlowID != "" && e.FlowID != f.FlowID {
		return false
	}
	if f.Status != "" && e.Status != f.Status {
		return false
	}
	return true
}
