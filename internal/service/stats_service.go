package service

import (
	"sort"

	"orchestr/internal/model"
)

type StatsOverview struct {
	TotalExecutions     int     `json:"total_executions"`
	CompletedExecutions int     `json:"completed_executions"`
	FailedExecutions    int     `json:"failed_executions"`
	TimeoutExecutions   int     `json:"timeout_executions"`
	SuccessRate         float64 `json:"success_rate"`
	AvgDurationMs       int     `json:"avg_duration_ms"`
}

type FlowDistribution struct {
	FlowID string `json:"flow_id"`
	Count  int    `json:"count"`
}

type ServiceDistribution struct {
	ServiceID string `json:"service_id"`
	Count     int    `json:"count"`
}

type FailureReasonDistribution struct {
	Reason string `json:"reason"`
	Count  int    `json:"count"`
}

func (s *Service) StatsOverview() (*StatsOverview, error) {
	all := s.store.ListExecutions()
	var completed, failed, timeout int
	var totalDuration int
	var durationCount int
	for _, e := range all {
		if e.Status == model.ExecutionCompleted {
			completed++
		}
		if e.Status == model.ExecutionFailed {
			failed++
		}
		if e.Status == model.ExecutionTimeout {
			timeout++
		}
		stepExecs := s.store.ListStepExecutionsByExecutionID(e.ID)
		for _, se := range stepExecs {
			totalDuration += se.DurationMs
			durationCount++
		}
	}
	total := len(all)
	var rate float64
	if total > 0 {
		rate = float64(completed) / float64(total) * 100
	}
	var avg int
	if durationCount > 0 {
		avg = totalDuration / durationCount
	}
	return &StatsOverview{
		TotalExecutions:     total,
		CompletedExecutions: completed,
		FailedExecutions:    failed,
		TimeoutExecutions:   timeout,
		SuccessRate:         rate,
		AvgDurationMs:       avg,
	}, nil
}

func (s *Service) StatsSuccessRate() map[string]float64 {
	all := s.store.ListExecutions()
	flowTotal := make(map[string]int)
	flowSuccess := make(map[string]int)
	for _, e := range all {
		flowTotal[e.FlowID]++
		if e.Status == model.ExecutionCompleted {
			flowSuccess[e.FlowID]++
		}
	}
	result := make(map[string]float64)
	for fid, total := range flowTotal {
		if total > 0 {
			result[fid] = float64(flowSuccess[fid]) / float64(total) * 100
		}
	}
	return result
}

func (s *Service) StatsByFlow() []FlowDistribution {
	all := s.store.ListExecutions()
	m := make(map[string]int)
	for _, e := range all {
		m[e.FlowID]++
	}
	result := make([]FlowDistribution, 0, len(m))
	for fid, c := range m {
		result = append(result, FlowDistribution{FlowID: fid, Count: c})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})
	return result
}

func (s *Service) StatsByService() []ServiceDistribution {
	all := s.store.ListStepExecutions()
	m := make(map[string]int)
	for _, se := range all {
		step, err := s.store.GetStep(se.StepID)
		if err != nil {
			continue
		}
		m[step.ServiceID]++
	}
	result := make([]ServiceDistribution, 0, len(m))
	for sid, c := range m {
		result = append(result, ServiceDistribution{ServiceID: sid, Count: c})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})
	return result
}

func (s *Service) StatsFailureReasons() []FailureReasonDistribution {
	all := s.store.ListExecutions()
	m := make(map[string]int)
	for _, e := range all {
		if e.Status == model.ExecutionFailed && e.ErrorMsg != "" {
			m[e.ErrorMsg]++
		}
	}
	result := make([]FailureReasonDistribution, 0, len(m))
	for r, c := range m {
		result = append(result, FailureReasonDistribution{Reason: r, Count: c})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})
	return result
}
