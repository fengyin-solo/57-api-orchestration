package service

import (
	"encoding/json"
	"testing"

	"orchestr/internal/config"
	"orchestr/internal/model"
	"orchestr/internal/store"
	"orchestr/pkg/logger"
)

func newTestService() *Service {
	cfg := &config.Config{MaxPageSize: 100}
	log := logger.NewLevel(logger.LevelError)
	return New(store.NewMemoryStore(), log, cfg)
}

func TestCreateService(t *testing.T) {
	s := newTestService()
	sv, err := s.CreateService(model.Service{Name: "svc", BaseURL: "http://a", AuthType: model.ServiceAuthNone, Status: model.ServiceActive})
	if err != nil {
		t.Fatalf("create service: %v", err)
	}
	if sv.ID == "" {
		t.Fatal("service id empty")
	}
	_, err = s.CreateService(model.Service{Name: "svc", BaseURL: "http://b", AuthType: model.ServiceAuthNone, Status: model.ServiceActive})
	if err == nil {
		t.Fatal("expected conflict")
	}
}

func TestServiceListPagination(t *testing.T) {
	s := newTestService()
	for i := 0; i < 5; i++ {
		_, _ = s.CreateService(model.Service{Name: "svc" + string(rune('0'+i)), BaseURL: "http://a", AuthType: model.ServiceAuthNone, Status: model.ServiceActive})
	}
	items, total, err := s.ListServices(model.ServiceFilter{}, 1, 2)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total != 5 {
		t.Fatalf("expected total 5, got %d", total)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
}

func TestFlowCRUD(t *testing.T) {
	s := newTestService()
	f, err := s.CreateFlow(model.Flow{Name: "flow1", Status: model.FlowActive})
	if err != nil {
		t.Fatalf("create flow: %v", err)
	}
	got, err := s.GetFlow(f.ID)
	if err != nil || got.Name != "flow1" {
		t.Fatalf("get flow: %v", err)
	}
	_, err = s.UpdateFlow(f.ID, model.Flow{Name: "flow1-updated", Status: model.FlowActive})
	if err != nil {
		t.Fatalf("update flow: %v", err)
	}
	if err := s.DeleteFlow(f.ID); err != nil {
		t.Fatalf("delete flow: %v", err)
	}
	if _, err := s.GetFlow(f.ID); err == nil {
		t.Fatal("expected not found")
	}
}

func TestStepWithForeignKey(t *testing.T) {
	s := newTestService()
	_, err := s.CreateStep(model.Step{FlowID: "nonexist", ServiceID: "nonexist", Method: "GET", Path: "/", Status: model.StepActive})
	if err == nil {
		t.Fatal("expected validation error for missing flow")
	}
	f, _ := s.CreateFlow(model.Flow{Name: "flow1", Status: model.FlowActive})
	sv, _ := s.CreateService(model.Service{Name: "svc1", BaseURL: "http://a", AuthType: model.ServiceAuthNone, Status: model.ServiceActive})
	st, err := s.CreateStep(model.Step{FlowID: f.ID, ServiceID: sv.ID, Method: "GET", Path: "/", Status: model.StepActive})
	if err != nil {
		t.Fatalf("create step: %v", err)
	}
	if st.FlowID != f.ID {
		t.Fatal("flow id mismatch")
	}
}

func TestExecutionStateMachine(t *testing.T) {
	s := newTestService()
	f, _ := s.CreateFlow(model.Flow{Name: "flow1", Status: model.FlowActive})
	e, err := s.CreateExecution(model.Execution{FlowID: f.ID})
	if err != nil {
		t.Fatalf("create execution: %v", err)
	}
	if e.Status != model.ExecutionPending {
		t.Fatalf("expected pending, got %s", e.Status)
	}
	e, err = s.RunExecution(e.ID)
	if err != nil {
		t.Fatalf("run execution: %v", err)
	}
	if e.Status != model.ExecutionRunning {
		t.Fatalf("expected running, got %s", e.Status)
	}
	_, err = s.RunExecution(e.ID)
	if err == nil {
		t.Fatal("expected error on double run")
	}
	e, err = s.CompleteExecution(e.ID, json.RawMessage(`{"ok":true}`))
	if err != nil {
		t.Fatalf("complete execution: %v", err)
	}
	if e.Status != model.ExecutionCompleted {
		t.Fatalf("expected completed, got %s", e.Status)
	}
	_, err = s.FailExecution(e.ID, "err")
	if err == nil {
		t.Fatal("expected error on fail after complete")
	}
}

func TestFailAndTimeoutExecution(t *testing.T) {
	s := newTestService()
	f, _ := s.CreateFlow(model.Flow{Name: "flow1", Status: model.FlowActive})
	e, _ := s.CreateExecution(model.Execution{FlowID: f.ID})
	e, _ = s.RunExecution(e.ID)
	e, err := s.FailExecution(e.ID, "something wrong")
	if err != nil {
		t.Fatalf("fail execution: %v", err)
	}
	if e.Status != model.ExecutionFailed {
		t.Fatalf("expected failed, got %s", e.Status)
	}

	e2, _ := s.CreateExecution(model.Execution{FlowID: f.ID})
	e2, _ = s.RunExecution(e2.ID)
	e2, err = s.TimeoutExecution(e2.ID)
	if err != nil {
		t.Fatalf("timeout execution: %v", err)
	}
	if e2.Status != model.ExecutionTimeout {
		t.Fatalf("expected timeout, got %s", e2.Status)
	}
}

func TestStepExecutionStateMachine(t *testing.T) {
	s := newTestService()
	f, _ := s.CreateFlow(model.Flow{Name: "flow1", Status: model.FlowActive})
	e, _ := s.CreateExecution(model.Execution{FlowID: f.ID})
	sv, _ := s.CreateService(model.Service{Name: "svc1", BaseURL: "http://a", AuthType: model.ServiceAuthNone, Status: model.ServiceActive})
	st, _ := s.CreateStep(model.Step{FlowID: f.ID, ServiceID: sv.ID, Method: "GET", Path: "/", Status: model.StepActive})
	se, err := s.CreateStepExecution(model.StepExecution{ExecutionID: e.ID, StepID: st.ID, Status: model.StepExecutionPending})
	if err != nil {
		t.Fatalf("create step execution: %v", err)
	}
	se, err = s.TransitionStepExecution(se.ID, model.StepExecutionRunning)
	if err != nil {
		t.Fatalf("transition to running: %v", err)
	}
	if se.Status != model.StepExecutionRunning {
		t.Fatalf("expected running, got %s", se.Status)
	}
	se, err = s.TransitionStepExecution(se.ID, model.StepExecutionCompleted)
	if err != nil {
		t.Fatalf("transition to completed: %v", err)
	}
	if se.Status != model.StepExecutionCompleted {
		t.Fatalf("expected completed, got %s", se.Status)
	}
	_, err = s.TransitionStepExecution(se.ID, model.StepExecutionFailed)
	if err == nil {
		t.Fatal("expected error on invalid transition")
	}
}

func TestTopoSort(t *testing.T) {
	steps := []*model.Step{
		{ID: "a"},
		{ID: "b"},
		{ID: "c"},
	}
	deps := []*model.Dependency{
		{StepID: "b", DependsOnStepID: "a"},
		{StepID: "c", DependsOnStepID: "b"},
	}
	order, err := TopoSort(steps, deps)
	if err != nil {
		t.Fatalf("topo sort: %v", err)
	}
	if len(order) != 3 {
		t.Fatalf("expected 3, got %d", len(order))
	}
	idx := make(map[string]int)
	for i, id := range order {
		idx[id] = i
	}
	if idx["a"] >= idx["b"] || idx["b"] >= idx["c"] {
		t.Fatal("topo order invalid")
	}
}

func TestTopoSortCycle(t *testing.T) {
	steps := []*model.Step{{ID: "a"}, {ID: "b"}}
	deps := []*model.Dependency{
		{StepID: "b", DependsOnStepID: "a"},
		{StepID: "a", DependsOnStepID: "b"},
	}
	_, err := TopoSort(steps, deps)
	if err == nil {
		t.Fatal("expected cycle error")
	}
}

func TestDependencySatisfied(t *testing.T) {
	deps := []*model.Dependency{
		{StepID: "b", DependsOnStepID: "a"},
	}
	graph := BuildDependencyGraph(deps)
	completed := map[string]bool{"a": true}
	if !AreDependenciesSatisfied("b", graph, completed) {
		t.Fatal("expected satisfied")
	}
	if !AreDependenciesSatisfied("a", graph, completed) {
		t.Fatal("expected satisfied for a")
	}
}

func TestExecuteFlow(t *testing.T) {
	s := newTestService()
	f, _ := s.CreateFlow(model.Flow{Name: "flow1", Status: model.FlowActive})
	sv, _ := s.CreateService(model.Service{Name: "svc1", BaseURL: "http://a", AuthType: model.ServiceAuthNone, Status: model.ServiceActive})
	st1, _ := s.CreateStep(model.Step{FlowID: f.ID, ServiceID: sv.ID, Method: "GET", Path: "/a", Order: 1, Status: model.StepActive})
	st2, _ := s.CreateStep(model.Step{FlowID: f.ID, ServiceID: sv.ID, Method: "GET", Path: "/b", Order: 2, Status: model.StepActive})
	_, _ = s.CreateDependency(model.Dependency{FlowID: f.ID, StepID: st2.ID, DependsOnStepID: st1.ID})

	exec, err := s.ExecuteFlow(f.ID, json.RawMessage(`{"key":"val"}`))
	if err != nil {
		t.Fatalf("execute flow: %v", err)
	}
	if exec.Status != model.ExecutionCompleted {
		t.Fatalf("expected completed, got %s", exec.Status)
	}
}

func TestExecuteFlowWithCondition(t *testing.T) {
	s := newTestService()
	f, _ := s.CreateFlow(model.Flow{Name: "flow1", Status: model.FlowActive})
	sv, _ := s.CreateService(model.Service{Name: "svc1", BaseURL: "http://a", AuthType: model.ServiceAuthNone, Status: model.ServiceActive})
	st1, _ := s.CreateStep(model.Step{FlowID: f.ID, ServiceID: sv.ID, Method: "GET", Path: "/a", Order: 1, Status: model.StepActive})
	st2, _ := s.CreateStep(model.Step{FlowID: f.ID, ServiceID: sv.ID, Method: "GET", Path: "/b", Order: 2, Status: model.StepActive})
	_, _ = s.CreateCondition(model.Condition{FlowID: f.ID, StepID: st1.ID, Expression: "true", TrueStepID: st2.ID, Status: model.ConditionActive})

	exec, err := s.ExecuteFlow(f.ID, json.RawMessage(`{"result":true}`))
	if err != nil {
		t.Fatalf("execute flow with condition: %v", err)
	}
	if exec.Status != model.ExecutionCompleted {
		t.Fatalf("expected completed, got %s", exec.Status)
	}
}

func TestExecuteFlowWithParallel(t *testing.T) {
	s := newTestService()
	f, _ := s.CreateFlow(model.Flow{Name: "flow1", Status: model.FlowActive})
	sv, _ := s.CreateService(model.Service{Name: "svc1", BaseURL: "http://a", AuthType: model.ServiceAuthNone, Status: model.ServiceActive})
	st1, _ := s.CreateStep(model.Step{FlowID: f.ID, ServiceID: sv.ID, Method: "GET", Path: "/a", Order: 1, Status: model.StepActive})
	st2, _ := s.CreateStep(model.Step{FlowID: f.ID, ServiceID: sv.ID, Method: "GET", Path: "/b", Order: 2, Status: model.StepActive})
	_, _ = s.CreateParallelConfig(model.ParallelConfig{FlowID: f.ID, StepIDs: []string{st1.ID, st2.ID}, Status: model.ParallelConfigActive})

	exec, err := s.ExecuteFlow(f.ID, json.RawMessage(`{}`))
	if err != nil {
		t.Fatalf("execute flow with parallel: %v", err)
	}
	if exec.Status != model.ExecutionCompleted {
		t.Fatalf("expected completed, got %s", exec.Status)
	}
}

func TestBatchCreateSteps(t *testing.T) {
	s := newTestService()
	f, _ := s.CreateFlow(model.Flow{Name: "flow1", Status: model.FlowActive})
	sv, _ := s.CreateService(model.Service{Name: "svc1", BaseURL: "http://a", AuthType: model.ServiceAuthNone, Status: model.ServiceActive})
	inputs := []model.Step{
		{FlowID: f.ID, ServiceID: sv.ID, Method: "GET", Path: "/1", Status: model.StepActive},
		{FlowID: f.ID, ServiceID: sv.ID, Method: "GET", Path: "/2", Status: model.StepActive},
	}
	items, err := s.BatchCreateSteps(inputs)
	if err != nil {
		t.Fatalf("batch create steps: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2, got %d", len(items))
	}
}

func TestBatchExecute(t *testing.T) {
	s := newTestService()
	f, _ := s.CreateFlow(model.Flow{Name: "flow1", Status: model.FlowActive})
	sv, _ := s.CreateService(model.Service{Name: "svc1", BaseURL: "http://a", AuthType: model.ServiceAuthNone, Status: model.ServiceActive})
	_, _ = s.CreateStep(model.Step{FlowID: f.ID, ServiceID: sv.ID, Method: "GET", Path: "/", Order: 1, Status: model.StepActive})
	inputs := []json.RawMessage{json.RawMessage(`{}`), json.RawMessage(`{}`)}
	items, err := s.BatchExecute(f.ID, inputs)
	if err != nil {
		t.Fatalf("batch execute: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2, got %d", len(items))
	}
}

func TestRetryPolicyCalculateDelay(t *testing.T) {
	rp := &model.RetryPolicy{InitialDelayMs: 100, BackoffType: model.BackoffFixed}
	if rp.CalculateDelay(0) != 100 {
		t.Fatalf("fixed delay mismatch")
	}
	rp.BackoffType = model.BackoffLinear
	if rp.CalculateDelay(1) != 200 {
		t.Fatalf("linear delay mismatch")
	}
	rp.BackoffType = model.BackoffExponential
	if rp.CalculateDelay(2) != 400 {
		t.Fatalf("exponential delay mismatch")
	}
}

func TestCircuitBreakerValidation(t *testing.T) {
	s := newTestService()
	_, err := s.CreateCircuitBreaker(model.CircuitBreaker{ServiceID: "nonexist", Threshold: 5, CooldownMs: 1000, Status: model.CircuitBreakerActive})
	if err == nil {
		t.Fatal("expected validation error for missing service")
	}
}

func TestStats(t *testing.T) {
	s := newTestService()
	f, _ := s.CreateFlow(model.Flow{Name: "flow1", Status: model.FlowActive})
	sv, _ := s.CreateService(model.Service{Name: "svc1", BaseURL: "http://a", AuthType: model.ServiceAuthNone, Status: model.ServiceActive})
	_, _ = s.CreateStep(model.Step{FlowID: f.ID, ServiceID: sv.ID, Method: "GET", Path: "/", Order: 1, Status: model.StepActive})
	_, _ = s.ExecuteFlow(f.ID, json.RawMessage(`{}`))
	_, _ = s.ExecuteFlow(f.ID, json.RawMessage(`{}`))

	overview, err := s.StatsOverview()
	if err != nil {
		t.Fatalf("stats overview: %v", err)
	}
	if overview.TotalExecutions != 2 {
		t.Fatalf("expected total 2, got %d", overview.TotalExecutions)
	}
	if overview.SuccessRate != 100 {
		t.Fatalf("expected 100%% success, got %f", overview.SuccessRate)
	}
	if len(s.StatsByFlow()) != 1 {
		t.Fatalf("expected 1 flow distribution")
	}
	if len(s.StatsByService()) != 1 {
		t.Fatalf("expected 1 service distribution")
	}
	if len(s.StatsFailureReasons()) != 0 {
		t.Fatalf("expected 0 failure reasons")
	}
}

func TestConditionEvaluate(t *testing.T) {
	c := &model.Condition{Expression: "result"}
	if !c.Evaluate(map[string]interface{}{"result": true}) {
		t.Fatal("expected true")
	}
	if c.Evaluate(map[string]interface{}{"result": false}) {
		t.Fatal("expected false")
	}
	if !c.Evaluate(map[string]interface{}{}) {
		t.Fatal("expected default true")
	}
}

func TestExportAll(t *testing.T) {
	s := newTestService()
	f, _ := s.CreateFlow(model.Flow{Name: "flow1", Status: model.FlowActive})
	data := s.ExportAll()
	if data == nil {
		t.Fatal("export nil")
	}
	flows, ok := data["flows"].(map[string]*model.Flow)
	if !ok || len(flows) != 1 {
		t.Fatalf("expected 1 flow in export, got %d", len(flows))
	}
	if flows[f.ID] == nil {
		t.Fatal("exported flow missing")
	}
}

func TestEnvVarCRUD(t *testing.T) {
	s := newTestService()
	f, _ := s.CreateFlow(model.Flow{Name: "flow1", Status: model.FlowActive})
	e, err := s.CreateEnvVar(model.EnvVar{FlowID: f.ID, Key: "k1", Value: "v1", Status: model.EnvVarActive})
	if err != nil {
		t.Fatalf("create env var: %v", err)
	}
	if e.Key != "k1" {
		t.Fatal("key mismatch")
	}
	m := s.BuildEnvMap(f.ID)
	if m["k1"] != "v1" {
		t.Fatalf("env map mismatch")
	}
	_, err = s.UpdateEnvVar(e.ID, model.EnvVar{FlowID: f.ID, Key: "k1", Value: "v2", Status: model.EnvVarActive})
	if err != nil {
		t.Fatalf("update env var: %v", err)
	}
	if err := s.DeleteEnvVar(e.ID); err != nil {
		t.Fatalf("delete env var: %v", err)
	}
}

func TestWebhookCRUD(t *testing.T) {
	s := newTestService()
	f, _ := s.CreateFlow(model.Flow{Name: "flow1", Status: model.FlowActive})
	w, err := s.CreateWebhook(model.Webhook{FlowID: f.ID, URL: "http://a", Method: "POST", Status: model.WebhookActive})
	if err != nil {
		t.Fatalf("create webhook: %v", err)
	}
	if w.URL != "http://a" {
		t.Fatal("url mismatch")
	}
	if err := s.TriggerWebhooks(f.ID, "{}"); err != nil {
		t.Fatalf("trigger webhooks: %v", err)
	}
	_, err = s.UpdateWebhook(w.ID, model.Webhook{FlowID: f.ID, URL: "http://b", Method: "POST", Status: model.WebhookActive})
	if err != nil {
		t.Fatalf("update webhook: %v", err)
	}
	if err := s.DeleteWebhook(w.ID); err != nil {
		t.Fatalf("delete webhook: %v", err)
	}
}

func TestTemplateCRUD(t *testing.T) {
	s := newTestService()
	tpl, err := s.CreateTemplate(model.Template{Name: "tpl1", Content: "hello {{name}}", ContentType: "text/plain", Status: model.TemplateActive})
	if err != nil {
		t.Fatalf("create template: %v", err)
	}
	if tpl.Name != "tpl1" {
		t.Fatal("name mismatch")
	}
	result, err := s.RenderTemplate("tpl1", map[string]string{"name": "world"})
	if err != nil {
		t.Fatalf("render template: %v", err)
	}
	if result != "hello world" {
		t.Fatalf("render result mismatch: %s", result)
	}
	_, err = s.UpdateTemplate(tpl.ID, model.Template{Name: "tpl1", Content: "hi {{name}}", ContentType: "text/plain", Status: model.TemplateActive})
	if err != nil {
		t.Fatalf("update template: %v", err)
	}
	if err := s.DeleteTemplate(tpl.ID); err != nil {
		t.Fatalf("delete template: %v", err)
	}
}

func TestWebhookLogCRUD(t *testing.T) {
	s := newTestService()
	f, _ := s.CreateFlow(model.Flow{Name: "flow1", Status: model.FlowActive})
	w, _ := s.CreateWebhook(model.Webhook{FlowID: f.ID, URL: "http://a", Method: "POST", Status: model.WebhookActive})
	wl, err := s.CreateWebhookLog(model.WebhookLog{WebhookID: w.ID, Payload: "{}", Status: model.WebhookLogPending})
	if err != nil {
		t.Fatalf("create webhook log: %v", err)
	}
	if wl.WebhookID != w.ID {
		t.Fatal("webhook id mismatch")
	}
	_, err = s.UpdateWebhookLog(wl.ID, model.WebhookLog{WebhookID: w.ID, Payload: "{}", Status: model.WebhookLogDelivered})
	if err != nil {
		t.Fatalf("update webhook log: %v", err)
	}
	if err := s.DeleteWebhookLog(wl.ID); err != nil {
		t.Fatalf("delete webhook log: %v", err)
	}
}
