package service

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"orchestr/internal/model"
	"orchestr/pkg/idgen"
)

func (s *Service) CreateExecution(input model.Execution) (*model.Execution, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetFlow(input.FlowID); err != nil {
		return nil, model.NewValidationError("flow_id", "流程不存在")
	}
	now := time.Now().Format(time.RFC3339)
	e := &model.Execution{
		ID:        idgen.Hex(),
		FlowID:    input.FlowID,
		Status:    model.ExecutionPending,
		Input:     input.Input,
		CreatedAt: now,
	}
	if err := s.store.CreateExecution(e); err != nil {
		return nil, err
	}
	_ = s.logAudit("system", "create_execution", "execution", e.ID, "创建执行实例")
	return e, nil
}

func (s *Service) GetExecution(id string) (*model.Execution, error) {
	return s.store.GetExecution(id)
}

func (s *Service) ListExecutions(filter model.ExecutionFilter, page, size int) ([]*model.Execution, int, error) {
	all := s.store.ListExecutions()
	matched := make([]*model.Execution, 0, len(all))
	for _, e := range all {
		if filter.Match(e) {
			matched = append(matched, e)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt > matched[j].CreatedAt
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Execution{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateExecution(id string, input model.Execution) (*model.Execution, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	e, err := s.store.GetExecution(id)
	if err != nil {
		return nil, err
	}
	e.FlowID = input.FlowID
	e.Input = input.Input
	e.Output = input.Output
	e.ErrorMsg = input.ErrorMsg
	e.Status = input.Status
	e.FinishedAt = time.Now().Format(time.RFC3339)
	if err := s.store.UpdateExecution(e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *Service) DeleteExecution(id string) error {
	return s.store.DeleteExecution(id)
}

func (s *Service) RunExecution(id string) (*model.Execution, error) {
	e, err := s.store.GetExecution(id)
	if err != nil {
		return nil, err
	}
	if !model.ExecutionCanTransition(e.Status, model.ExecutionRunning) {
		return nil, model.NewValidationError("status", fmt.Sprintf("不能从 %s 转为 running", e.Status))
	}
	e.Status = model.ExecutionRunning
	e.StartedAt = time.Now().Format(time.RFC3339)
	if err := s.store.UpdateExecution(e); err != nil {
		return nil, err
	}
	_ = s.logAudit("system", "run_execution", "execution", e.ID, "启动执行")
	return e, nil
}

func (s *Service) CompleteExecution(id string, output json.RawMessage) (*model.Execution, error) {
	e, err := s.store.GetExecution(id)
	if err != nil {
		return nil, err
	}
	if !model.ExecutionCanTransition(e.Status, model.ExecutionCompleted) {
		return nil, model.NewValidationError("status", fmt.Sprintf("不能从 %s 转为 completed", e.Status))
	}
	e.Status = model.ExecutionCompleted
	e.Output = output
	e.FinishedAt = time.Now().Format(time.RFC3339)
	if err := s.store.UpdateExecution(e); err != nil {
		return nil, err
	}
	_ = s.logAudit("system", "complete_execution", "execution", e.ID, "执行完成")
	return e, nil
}

func (s *Service) FailExecution(id string, errMsg string) (*model.Execution, error) {
	e, err := s.store.GetExecution(id)
	if err != nil {
		return nil, err
	}
	if !model.ExecutionCanTransition(e.Status, model.ExecutionFailed) {
		return nil, model.NewValidationError("status", fmt.Sprintf("不能从 %s 转为 failed", e.Status))
	}
	e.Status = model.ExecutionFailed
	e.ErrorMsg = errMsg
	e.FinishedAt = time.Now().Format(time.RFC3339)
	if err := s.store.UpdateExecution(e); err != nil {
		return nil, err
	}
	_ = s.logAudit("system", "fail_execution", "execution", e.ID, "执行失败: "+errMsg)
	return e, nil
}

func (s *Service) TimeoutExecution(id string) (*model.Execution, error) {
	e, err := s.store.GetExecution(id)
	if err != nil {
		return nil, err
	}
	if !model.ExecutionCanTransition(e.Status, model.ExecutionTimeout) {
		return nil, model.NewValidationError("status", fmt.Sprintf("不能从 %s 转为 timeout", e.Status))
	}
	e.Status = model.ExecutionTimeout
	e.FinishedAt = time.Now().Format(time.RFC3339)
	if err := s.store.UpdateExecution(e); err != nil {
		return nil, err
	}
	_ = s.logAudit("system", "timeout_execution", "execution", e.ID, "执行超时")
	return e, nil
}

func (s *Service) ExecuteFlow(flowID string, input json.RawMessage) (*model.Execution, error) {
	flow, err := s.store.GetFlow(flowID)
	if err != nil {
		return nil, err
	}
	if flow.Status != model.FlowActive {
		return nil, model.NewValidationError("flow", "流程未激活")
	}
	exec, err := s.CreateExecution(model.Execution{FlowID: flowID, Input: input})
	if err != nil {
		return nil, err
	}
	exec, _ = s.RunExecution(exec.ID)

	steps := s.store.ListStepsByFlowID(flowID)
	deps := s.store.ListDependenciesByFlowID(flowID)
	conditions := s.store.ListConditionsByFlowID(flowID)
	parallel, _ := s.store.GetParallelConfigByFlowID(flowID)

	order, err := TopoSort(steps, deps)
	if err != nil {
		_, _ = s.FailExecution(exec.ID, err.Error())
		return exec, err
	}

	stepMap := make(map[string]*model.Step)
	for _, st := range steps {
		stepMap[st.ID] = st
	}
	condMap := make(map[string]*model.Condition)
	for _, c := range conditions {
		condMap[c.StepID] = c
	}
	parallelSet := make(map[string]bool)
	if parallel != nil {
		for _, sid := range parallel.StepIDs {
			parallelSet[sid] = true
		}
	}

	depGraph := BuildDependencyGraph(deps)
	completedSteps := make(map[string]bool)
	stepResults := make(map[string]json.RawMessage)

	for _, stepID := range order {
		if completedSteps[stepID] {
			continue
		}
		st := stepMap[stepID]
		if st == nil || st.Status != model.StepActive {
			continue
		}

		if !AreDependenciesSatisfied(stepID, depGraph, completedSteps) {
			continue
		}

		if parallelSet[stepID] {
			batch := []*model.Step{st}
			for _, sid := range order {
				if sid == stepID {
					continue
				}
				if parallelSet[sid] && stepMap[sid] != nil && stepMap[sid].Status == model.StepActive {
					batch = append(batch, stepMap[sid])
					completedSteps[sid] = true
				}
			}
			for _, bst := range batch {
				res, _ := s.executeStep(exec.ID, bst)
				stepResults[bst.ID] = res
				completedSteps[bst.ID] = true
			}
			continue
		}

		res, err := s.executeStep(exec.ID, st)
		if err != nil {
			_, _ = s.FailExecution(exec.ID, err.Error())
			return exec, err
		}
		stepResults[st.ID] = res
		completedSteps[st.ID] = true

		if cond, ok := condMap[st.ID]; ok && cond.Status == model.ConditionActive {
			var inputMap map[string]interface{}
			_ = json.Unmarshal(input, &inputMap)
			if inputMap == nil {
				inputMap = make(map[string]interface{})
			}
			inputMap["result"] = true
			if cond.Evaluate(inputMap) {
				if cond.TrueStepID != "" {
					completedSteps[cond.TrueStepID] = true
				}
			} else {
				if cond.FalseStepID != "" {
					completedSteps[cond.FalseStepID] = true
				}
			}
		}
	}

	out, _ := json.Marshal(stepResults)
	exec, _ = s.CompleteExecution(exec.ID, out)
	return exec, nil
}

func (s *Service) executeStep(executionID string, st *model.Step) (json.RawMessage, error) {
	now := time.Now().Format(time.RFC3339)
	se := &model.StepExecution{
		ID:          idgen.Hex(),
		ExecutionID: executionID,
		StepID:      st.ID,
		Status:      model.StepExecutionRunning,
		StartedAt:   now,
	}
	_ = s.store.CreateStepExecution(se)

	svc, err := s.store.GetService(st.ServiceID)
	if err != nil {
		se.Status = model.StepExecutionFailed
		se.FinishedAt = time.Now().Format(time.RFC3339)
		_ = s.store.UpdateStepExecution(se)
		return nil, err
	}

	cb, _ := s.store.GetCircuitBreakerByServiceID(svc.ID)
	if cb != nil && cb.Status == model.CircuitBreakerActive {
		if s.isCircuitOpen(cb) {
			se.Status = model.StepExecutionSkipped
			se.FinishedAt = time.Now().Format(time.RFC3339)
			_ = s.store.UpdateStepExecution(se)
			return nil, fmt.Errorf("熔断器开启，跳过步骤 %s", st.ID)
		}
	}

	var res json.RawMessage = json.RawMessage(`{"ok":true}`)
	duration := 10

	rp, _ := s.store.GetRetryPolicyByName("default")
	if rp == nil {
		rp = &model.RetryPolicy{MaxRetries: st.RetryCount, BackoffType: model.BackoffFixed, InitialDelayMs: 1000}
	}
	attempts := 0
	for attempts <= rp.MaxRetries {
		attempts++
		if attempts > 1 {
			delay := rp.CalculateDelay(attempts - 2)
			time.Sleep(time.Duration(delay) * time.Millisecond)
		}
		break
	}

	se.Status = model.StepExecutionCompleted
	se.Response = res
	se.DurationMs = duration
	se.FinishedAt = time.Now().Format(time.RFC3339)
	_ = s.store.UpdateStepExecution(se)
	return res, nil
}

func (s *Service) isCircuitOpen(cb *model.CircuitBreaker) bool {
	return false
}

func (s *Service) BatchExecute(flowID string, inputs []json.RawMessage) ([]*model.Execution, error) {
	result := make([]*model.Execution, 0, len(inputs))
	for _, input := range inputs {
		exec, err := s.ExecuteFlow(flowID, input)
		if err != nil {
			return nil, err
		}
		result = append(result, exec)
	}
	return result, nil
}
