// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"orchestr/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法，便于测试时替换实现。
type Store interface {
	// Service
	CreateService(s *model.Service) error
	GetService(id string) (*model.Service, error)
	GetServiceByName(name string) (*model.Service, error)
	ListServices() []*model.Service
	UpdateService(s *model.Service) error
	DeleteService(id string) error

	// Flow
	CreateFlow(f *model.Flow) error
	GetFlow(id string) (*model.Flow, error)
	GetFlowByName(name string) (*model.Flow, error)
	ListFlows() []*model.Flow
	UpdateFlow(f *model.Flow) error
	DeleteFlow(id string) error

	// Step
	CreateStep(s *model.Step) error
	GetStep(id string) (*model.Step, error)
	ListSteps() []*model.Step
	UpdateStep(s *model.Step) error
	DeleteStep(id string) error
	ListStepsByFlowID(flowID string) []*model.Step

	// Execution
	CreateExecution(e *model.Execution) error
	GetExecution(id string) (*model.Execution, error)
	ListExecutions() []*model.Execution
	UpdateExecution(e *model.Execution) error
	DeleteExecution(id string) error

	// StepExecution
	CreateStepExecution(se *model.StepExecution) error
	GetStepExecution(id string) (*model.StepExecution, error)
	ListStepExecutions() []*model.StepExecution
	UpdateStepExecution(se *model.StepExecution) error
	DeleteStepExecution(id string) error
	ListStepExecutionsByExecutionID(executionID string) []*model.StepExecution

	// Dependency
	CreateDependency(d *model.Dependency) error
	GetDependency(id string) (*model.Dependency, error)
	ListDependencies() []*model.Dependency
	UpdateDependency(d *model.Dependency) error
	DeleteDependency(id string) error
	ListDependenciesByFlowID(flowID string) []*model.Dependency

	// RetryPolicy
	CreateRetryPolicy(r *model.RetryPolicy) error
	GetRetryPolicy(id string) (*model.RetryPolicy, error)
	GetRetryPolicyByName(name string) (*model.RetryPolicy, error)
	ListRetryPolicies() []*model.RetryPolicy
	UpdateRetryPolicy(r *model.RetryPolicy) error
	DeleteRetryPolicy(id string) error

	// CircuitBreaker
	CreateCircuitBreaker(c *model.CircuitBreaker) error
	GetCircuitBreaker(id string) (*model.CircuitBreaker, error)
	GetCircuitBreakerByServiceID(serviceID string) (*model.CircuitBreaker, error)
	ListCircuitBreakers() []*model.CircuitBreaker
	UpdateCircuitBreaker(c *model.CircuitBreaker) error
	DeleteCircuitBreaker(id string) error

	// ParallelConfig
	CreateParallelConfig(p *model.ParallelConfig) error
	GetParallelConfig(id string) (*model.ParallelConfig, error)
	ListParallelConfigs() []*model.ParallelConfig
	UpdateParallelConfig(p *model.ParallelConfig) error
	DeleteParallelConfig(id string) error
	GetParallelConfigByFlowID(flowID string) (*model.ParallelConfig, error)

	// Condition
	CreateCondition(c *model.Condition) error
	GetCondition(id string) (*model.Condition, error)
	ListConditions() []*model.Condition
	UpdateCondition(c *model.Condition) error
	DeleteCondition(id string) error
	ListConditionsByFlowID(flowID string) []*model.Condition

	// AuditLog
	CreateAuditLog(a *model.AuditLog) error
	GetAuditLog(id string) (*model.AuditLog, error)
	ListAuditLogs() []*model.AuditLog
	DeleteAuditLog(id string) error

	// EnvVar
	CreateEnvVar(e *model.EnvVar) error
	GetEnvVar(id string) (*model.EnvVar, error)
	ListEnvVars() []*model.EnvVar
	UpdateEnvVar(e *model.EnvVar) error
	DeleteEnvVar(id string) error
	ListEnvVarsByFlowID(flowID string) []*model.EnvVar

	// Webhook
	CreateWebhook(w *model.Webhook) error
	GetWebhook(id string) (*model.Webhook, error)
	ListWebhooks() []*model.Webhook
	UpdateWebhook(w *model.Webhook) error
	DeleteWebhook(id string) error
	ListWebhooksByFlowID(flowID string) []*model.Webhook

	// Template
	CreateTemplate(t *model.Template) error
	GetTemplate(id string) (*model.Template, error)
	GetTemplateByName(name string) (*model.Template, error)
	ListTemplates() []*model.Template
	UpdateTemplate(t *model.Template) error
	DeleteTemplate(id string) error

	// WebhookLog
	CreateWebhookLog(w *model.WebhookLog) error
	GetWebhookLog(id string) (*model.WebhookLog, error)
	ListWebhookLogs() []*model.WebhookLog
	UpdateWebhookLog(w *model.WebhookLog) error
	DeleteWebhookLog(id string) error
	ListWebhookLogsByWebhookID(webhookID string) []*model.WebhookLog

	// Snapshot export
	ExportAll() map[string]interface{}
}
