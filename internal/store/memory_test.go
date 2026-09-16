package store

import (
	"testing"

	"orchestr/internal/model"
)

func TestMemoryStoreService(t *testing.T) {
	s := NewMemoryStore()
	sv := &model.Service{ID: "s1", Name: "svc1", BaseURL: "http://a", AuthType: model.ServiceAuthNone, TimeoutMs: 1000, RetryCount: 0, Status: model.ServiceActive}
	if err := s.CreateService(sv); err != nil {
		t.Fatalf("create service: %v", err)
	}
	if err := s.CreateService(&model.Service{ID: "s2", Name: "svc1", BaseURL: "http://b", AuthType: model.ServiceAuthNone, Status: model.ServiceActive}); err != ErrConflict {
		t.Fatalf("expected conflict, got %v", err)
	}
	got, err := s.GetService("s1")
	if err != nil || got.Name != "svc1" {
		t.Fatalf("get service: %v", err)
	}
	if _, err := s.GetService("s9"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	if len(s.ListServices()) != 1 {
		t.Fatalf("expected 1 service, got %d", len(s.ListServices()))
	}
	sv.Name = "svc1-updated"
	if err := s.UpdateService(sv); err != nil {
		t.Fatalf("update service: %v", err)
	}
	if err := s.DeleteService("s1"); err != nil {
		t.Fatalf("delete service: %v", err)
	}
	if _, err := s.GetService("s1"); err != ErrNotFound {
		t.Fatalf("expected not found after delete, got %v", err)
	}
}

func TestMemoryStoreFlow(t *testing.T) {
	s := NewMemoryStore()
	f := &model.Flow{ID: "f1", Name: "flow1", Status: model.FlowActive}
	if err := s.CreateFlow(f); err != nil {
		t.Fatalf("create flow: %v", err)
	}
	if err := s.CreateFlow(&model.Flow{ID: "f2", Name: "flow1", Status: model.FlowActive}); err != ErrConflict {
		t.Fatalf("expected conflict, got %v", err)
	}
	got, err := s.GetFlow("f1")
	if err != nil || got.Name != "flow1" {
		t.Fatalf("get flow: %v", err)
	}
	if _, err := s.GetFlow("f9"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	if len(s.ListFlows()) != 1 {
		t.Fatalf("expected 1 flow, got %d", len(s.ListFlows()))
	}
	f.Name = "flow1-updated"
	if err := s.UpdateFlow(f); err != nil {
		t.Fatalf("update flow: %v", err)
	}
	if err := s.DeleteFlow("f1"); err != nil {
		t.Fatalf("delete flow: %v", err)
	}
	if _, err := s.GetFlow("f1"); err != ErrNotFound {
		t.Fatalf("expected not found after delete, got %v", err)
	}
}

func TestMemoryStoreStep(t *testing.T) {
	s := NewMemoryStore()
	st := &model.Step{ID: "st1", FlowID: "f1", ServiceID: "s1", Method: "GET", Path: "/", Order: 1, Status: model.StepActive}
	if err := s.CreateStep(st); err != nil {
		t.Fatalf("create step: %v", err)
	}
	got, err := s.GetStep("st1")
	if err != nil || got.Method != "GET" {
		t.Fatalf("get step: %v", err)
	}
	if _, err := s.GetStep("st9"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	if len(s.ListSteps()) != 1 {
		t.Fatalf("expected 1 step, got %d", len(s.ListSteps()))
	}
	st.Method = "POST"
	if err := s.UpdateStep(st); err != nil {
		t.Fatalf("update step: %v", err)
	}
	if err := s.DeleteStep("st1"); err != nil {
		t.Fatalf("delete step: %v", err)
	}
	if _, err := s.GetStep("st1"); err != ErrNotFound {
		t.Fatalf("expected not found after delete, got %v", err)
	}
}

func TestMemoryStoreExecution(t *testing.T) {
	s := NewMemoryStore()
	e := &model.Execution{ID: "e1", FlowID: "f1", Status: model.ExecutionPending}
	if err := s.CreateExecution(e); err != nil {
		t.Fatalf("create execution: %v", err)
	}
	got, err := s.GetExecution("e1")
	if err != nil || got.Status != model.ExecutionPending {
		t.Fatalf("get execution: %v", err)
	}
	if _, err := s.GetExecution("e9"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	if len(s.ListExecutions()) != 1 {
		t.Fatalf("expected 1 execution, got %d", len(s.ListExecutions()))
	}
	e.Status = model.ExecutionRunning
	if err := s.UpdateExecution(e); err != nil {
		t.Fatalf("update execution: %v", err)
	}
	if err := s.DeleteExecution("e1"); err != nil {
		t.Fatalf("delete execution: %v", err)
	}
	if _, err := s.GetExecution("e1"); err != ErrNotFound {
		t.Fatalf("expected not found after delete, got %v", err)
	}
}

func TestMemoryStoreStepExecution(t *testing.T) {
	s := NewMemoryStore()
	se := &model.StepExecution{ID: "se1", ExecutionID: "e1", StepID: "st1", Status: model.StepExecutionPending}
	if err := s.CreateStepExecution(se); err != nil {
		t.Fatalf("create step execution: %v", err)
	}
	got, err := s.GetStepExecution("se1")
	if err != nil || got.Status != model.StepExecutionPending {
		t.Fatalf("get step execution: %v", err)
	}
	if _, err := s.GetStepExecution("se9"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	if len(s.ListStepExecutions()) != 1 {
		t.Fatalf("expected 1 step execution, got %d", len(s.ListStepExecutions()))
	}
	se.Status = model.StepExecutionRunning
	if err := s.UpdateStepExecution(se); err != nil {
		t.Fatalf("update step execution: %v", err)
	}
	if err := s.DeleteStepExecution("se1"); err != nil {
		t.Fatalf("delete step execution: %v", err)
	}
	if _, err := s.GetStepExecution("se1"); err != ErrNotFound {
		t.Fatalf("expected not found after delete, got %v", err)
	}
}

func TestMemoryStoreDependency(t *testing.T) {
	s := NewMemoryStore()
	d := &model.Dependency{ID: "d1", FlowID: "f1", StepID: "st2", DependsOnStepID: "st1", Condition: ""}
	if err := s.CreateDependency(d); err != nil {
		t.Fatalf("create dependency: %v", err)
	}
	got, err := s.GetDependency("d1")
	if err != nil || got.StepID != "st2" {
		t.Fatalf("get dependency: %v", err)
	}
	if _, err := s.GetDependency("d9"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	if len(s.ListDependencies()) != 1 {
		t.Fatalf("expected 1 dependency, got %d", len(s.ListDependencies()))
	}
	d.Condition = "ok"
	if err := s.UpdateDependency(d); err != nil {
		t.Fatalf("update dependency: %v", err)
	}
	if err := s.DeleteDependency("d1"); err != nil {
		t.Fatalf("delete dependency: %v", err)
	}
	if _, err := s.GetDependency("d1"); err != ErrNotFound {
		t.Fatalf("expected not found after delete, got %v", err)
	}
}

func TestMemoryStoreRetryPolicy(t *testing.T) {
	s := NewMemoryStore()
	r := &model.RetryPolicy{ID: "r1", Name: "rp1", MaxRetries: 3, BackoffType: model.BackoffFixed, InitialDelayMs: 1000, Status: model.RetryPolicyActive}
	if err := s.CreateRetryPolicy(r); err != nil {
		t.Fatalf("create retry policy: %v", err)
	}
	if err := s.CreateRetryPolicy(&model.RetryPolicy{ID: "r2", Name: "rp1", MaxRetries: 3, BackoffType: model.BackoffFixed, InitialDelayMs: 1000, Status: model.RetryPolicyActive}); err != ErrConflict {
		t.Fatalf("expected conflict, got %v", err)
	}
	got, err := s.GetRetryPolicy("r1")
	if err != nil || got.Name != "rp1" {
		t.Fatalf("get retry policy: %v", err)
	}
	if _, err := s.GetRetryPolicy("r9"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	if len(s.ListRetryPolicies()) != 1 {
		t.Fatalf("expected 1 retry policy, got %d", len(s.ListRetryPolicies()))
	}
	r.MaxRetries = 5
	if err := s.UpdateRetryPolicy(r); err != nil {
		t.Fatalf("update retry policy: %v", err)
	}
	if err := s.DeleteRetryPolicy("r1"); err != nil {
		t.Fatalf("delete retry policy: %v", err)
	}
	if _, err := s.GetRetryPolicy("r1"); err != ErrNotFound {
		t.Fatalf("expected not found after delete, got %v", err)
	}
}

func TestMemoryStoreCircuitBreaker(t *testing.T) {
	s := NewMemoryStore()
	c := &model.CircuitBreaker{ID: "c1", ServiceID: "s1", Threshold: 5, CooldownMs: 1000, Status: model.CircuitBreakerActive}
	if err := s.CreateCircuitBreaker(c); err != nil {
		t.Fatalf("create circuit breaker: %v", err)
	}
	if err := s.CreateCircuitBreaker(&model.CircuitBreaker{ID: "c2", ServiceID: "s1", Threshold: 5, CooldownMs: 1000, Status: model.CircuitBreakerActive}); err != ErrConflict {
		t.Fatalf("expected conflict, got %v", err)
	}
	got, err := s.GetCircuitBreaker("c1")
	if err != nil || got.ServiceID != "s1" {
		t.Fatalf("get circuit breaker: %v", err)
	}
	if _, err := s.GetCircuitBreaker("c9"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	if len(s.ListCircuitBreakers()) != 1 {
		t.Fatalf("expected 1 circuit breaker, got %d", len(s.ListCircuitBreakers()))
	}
	c.Threshold = 10
	if err := s.UpdateCircuitBreaker(c); err != nil {
		t.Fatalf("update circuit breaker: %v", err)
	}
	if err := s.DeleteCircuitBreaker("c1"); err != nil {
		t.Fatalf("delete circuit breaker: %v", err)
	}
	if _, err := s.GetCircuitBreaker("c1"); err != ErrNotFound {
		t.Fatalf("expected not found after delete, got %v", err)
	}
}

func TestMemoryStoreParallelConfig(t *testing.T) {
	s := NewMemoryStore()
	p := &model.ParallelConfig{ID: "p1", FlowID: "f1", StepIDs: []string{"st1"}, Status: model.ParallelConfigActive}
	if err := s.CreateParallelConfig(p); err != nil {
		t.Fatalf("create parallel config: %v", err)
	}
	if err := s.CreateParallelConfig(&model.ParallelConfig{ID: "p2", FlowID: "f1", StepIDs: []string{"st2"}, Status: model.ParallelConfigActive}); err != ErrConflict {
		t.Fatalf("expected conflict, got %v", err)
	}
	got, err := s.GetParallelConfig("p1")
	if err != nil || got.FlowID != "f1" {
		t.Fatalf("get parallel config: %v", err)
	}
	if _, err := s.GetParallelConfig("p9"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	if len(s.ListParallelConfigs()) != 1 {
		t.Fatalf("expected 1 parallel config, got %d", len(s.ListParallelConfigs()))
	}
	p.StepIDs = []string{"st1", "st2"}
	if err := s.UpdateParallelConfig(p); err != nil {
		t.Fatalf("update parallel config: %v", err)
	}
	if err := s.DeleteParallelConfig("p1"); err != nil {
		t.Fatalf("delete parallel config: %v", err)
	}
	if _, err := s.GetParallelConfig("p1"); err != ErrNotFound {
		t.Fatalf("expected not found after delete, got %v", err)
	}
}

func TestMemoryStoreCondition(t *testing.T) {
	s := NewMemoryStore()
	c := &model.Condition{ID: "c1", FlowID: "f1", StepID: "st1", Expression: "true", TrueStepID: "st2", Status: model.ConditionActive}
	if err := s.CreateCondition(c); err != nil {
		t.Fatalf("create condition: %v", err)
	}
	got, err := s.GetCondition("c1")
	if err != nil || got.Expression != "true" {
		t.Fatalf("get condition: %v", err)
	}
	if _, err := s.GetCondition("c9"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	if len(s.ListConditions()) != 1 {
		t.Fatalf("expected 1 condition, got %d", len(s.ListConditions()))
	}
	c.Expression = "false"
	if err := s.UpdateCondition(c); err != nil {
		t.Fatalf("update condition: %v", err)
	}
	if err := s.DeleteCondition("c1"); err != nil {
		t.Fatalf("delete condition: %v", err)
	}
	if _, err := s.GetCondition("c1"); err != ErrNotFound {
		t.Fatalf("expected not found after delete, got %v", err)
	}
}

func TestMemoryStoreAuditLog(t *testing.T) {
	s := NewMemoryStore()
	a := &model.AuditLog{ID: "a1", Operator: "admin", Action: "create", TargetType: "flow", TargetID: "f1", Detail: "test"}
	if err := s.CreateAuditLog(a); err != nil {
		t.Fatalf("create audit log: %v", err)
	}
	got, err := s.GetAuditLog("a1")
	if err != nil || got.Operator != "admin" {
		t.Fatalf("get audit log: %v", err)
	}
	if _, err := s.GetAuditLog("a9"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	if len(s.ListAuditLogs()) != 1 {
		t.Fatalf("expected 1 audit log, got %d", len(s.ListAuditLogs()))
	}
	if err := s.DeleteAuditLog("a1"); err != nil {
		t.Fatalf("delete audit log: %v", err)
	}
	if _, err := s.GetAuditLog("a1"); err != ErrNotFound {
		t.Fatalf("expected not found after delete, got %v", err)
	}
}

func TestMemoryStoreEnvVar(t *testing.T) {
	s := NewMemoryStore()
	e := &model.EnvVar{ID: "e1", FlowID: "f1", Key: "K1", Value: "V1", Status: model.EnvVarActive}
	if err := s.CreateEnvVar(e); err != nil {
		t.Fatalf("create env var: %v", err)
	}
	got, err := s.GetEnvVar("e1")
	if err != nil || got.Key != "K1" {
		t.Fatalf("get env var: %v", err)
	}
	if _, err := s.GetEnvVar("e9"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	if len(s.ListEnvVars()) != 1 {
		t.Fatalf("expected 1 env var, got %d", len(s.ListEnvVars()))
	}
	e.Value = "V2"
	if err := s.UpdateEnvVar(e); err != nil {
		t.Fatalf("update env var: %v", err)
	}
	if err := s.DeleteEnvVar("e1"); err != nil {
		t.Fatalf("delete env var: %v", err)
	}
	if _, err := s.GetEnvVar("e1"); err != ErrNotFound {
		t.Fatalf("expected not found after delete, got %v", err)
	}
}

func TestMemoryStoreWebhook(t *testing.T) {
	s := NewMemoryStore()
	w := &model.Webhook{ID: "w1", FlowID: "f1", URL: "http://a", Method: "POST", Status: model.WebhookActive}
	if err := s.CreateWebhook(w); err != nil {
		t.Fatalf("create webhook: %v", err)
	}
	got, err := s.GetWebhook("w1")
	if err != nil || got.URL != "http://a" {
		t.Fatalf("get webhook: %v", err)
	}
	if _, err := s.GetWebhook("w9"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	if len(s.ListWebhooks()) != 1 {
		t.Fatalf("expected 1 webhook, got %d", len(s.ListWebhooks()))
	}
	w.URL = "http://b"
	if err := s.UpdateWebhook(w); err != nil {
		t.Fatalf("update webhook: %v", err)
	}
	if err := s.DeleteWebhook("w1"); err != nil {
		t.Fatalf("delete webhook: %v", err)
	}
	if _, err := s.GetWebhook("w1"); err != ErrNotFound {
		t.Fatalf("expected not found after delete, got %v", err)
	}
}

func TestMemoryStoreTemplate(t *testing.T) {
	s := NewMemoryStore()
	tpl := &model.Template{ID: "t1", Name: "tpl1", Content: "hello", Status: model.TemplateActive}
	if err := s.CreateTemplate(tpl); err != nil {
		t.Fatalf("create template: %v", err)
	}
	if err := s.CreateTemplate(&model.Template{ID: "t2", Name: "tpl1", Content: "world", Status: model.TemplateActive}); err != ErrConflict {
		t.Fatalf("expected conflict, got %v", err)
	}
	got, err := s.GetTemplate("t1")
	if err != nil || got.Name != "tpl1" {
		t.Fatalf("get template: %v", err)
	}
	if _, err := s.GetTemplate("t9"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	if len(s.ListTemplates()) != 1 {
		t.Fatalf("expected 1 template, got %d", len(s.ListTemplates()))
	}
	tpl.Content = "hello world"
	if err := s.UpdateTemplate(tpl); err != nil {
		t.Fatalf("update template: %v", err)
	}
	if err := s.DeleteTemplate("t1"); err != nil {
		t.Fatalf("delete template: %v", err)
	}
	if _, err := s.GetTemplate("t1"); err != ErrNotFound {
		t.Fatalf("expected not found after delete, got %v", err)
	}
}

func TestMemoryStoreWebhookLog(t *testing.T) {
	s := NewMemoryStore()
	w := &model.WebhookLog{ID: "wl1", WebhookID: "w1", Payload: "{}", Status: model.WebhookLogPending}
	if err := s.CreateWebhookLog(w); err != nil {
		t.Fatalf("create webhook log: %v", err)
	}
	got, err := s.GetWebhookLog("wl1")
	if err != nil || got.WebhookID != "w1" {
		t.Fatalf("get webhook log: %v", err)
	}
	if _, err := s.GetWebhookLog("wl9"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	if len(s.ListWebhookLogs()) != 1 {
		t.Fatalf("expected 1 webhook log, got %d", len(s.ListWebhookLogs()))
	}
	w.Status = model.WebhookLogDelivered
	if err := s.UpdateWebhookLog(w); err != nil {
		t.Fatalf("update webhook log: %v", err)
	}
	if err := s.DeleteWebhookLog("wl1"); err != nil {
		t.Fatalf("delete webhook log: %v", err)
	}
	if _, err := s.GetWebhookLog("wl1"); err != ErrNotFound {
		t.Fatalf("expected not found after delete, got %v", err)
	}
}
