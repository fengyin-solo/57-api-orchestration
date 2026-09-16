package model

const (
	ParallelConfigActive   = "active"
	ParallelConfigDisabled = "disabled"
)

type ParallelConfig struct {
	ID        string   `json:"id"`
	FlowID    string   `json:"flow_id"`
	StepIDs   []string `json:"step_ids"`
	Status    string   `json:"status"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

func (p *ParallelConfig) Validate() error {
	if p.FlowID == "" {
		return NewValidationError("flow_id", "流程ID不能为空")
	}
	if len(p.StepIDs) == 0 {
		return NewValidationError("step_ids", "并行步骤列表不能为空")
	}
	if p.Status == "" {
		p.Status = ParallelConfigActive
	}
	if p.Status != ParallelConfigActive && p.Status != ParallelConfigDisabled {
		return NewValidationError("status", "状态不合法")
	}
	return nil
}

type ParallelConfigFilter struct {
	FlowID string
	Status string
}

func (f ParallelConfigFilter) Match(p *ParallelConfig) bool {
	if f.FlowID != "" && p.FlowID != f.FlowID {
		return false
	}
	if f.Status != "" && p.Status != f.Status {
		return false
	}
	return true
}
