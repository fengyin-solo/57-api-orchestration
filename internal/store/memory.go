package store

import (
	"sync"

	"orchestr/internal/model"
)

type MemoryStore struct {
	mu              sync.RWMutex
	services        map[string]*model.Service
	flows           map[string]*model.Flow
	steps           map[string]*model.Step
	executions      map[string]*model.Execution
	stepExecutions  map[string]*model.StepExecution
	dependencies    map[string]*model.Dependency
	retryPolicies   map[string]*model.RetryPolicy
	circuitBreakers map[string]*model.CircuitBreaker
	parallelConfigs map[string]*model.ParallelConfig
	conditions      map[string]*model.Condition
	auditLogs       map[string]*model.AuditLog
	envVars         map[string]*model.EnvVar
	webhooks        map[string]*model.Webhook
	templates       map[string]*model.Template
	webhookLogs     map[string]*model.WebhookLog
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		services:        make(map[string]*model.Service),
		flows:           make(map[string]*model.Flow),
		steps:           make(map[string]*model.Step),
		executions:      make(map[string]*model.Execution),
		stepExecutions:  make(map[string]*model.StepExecution),
		dependencies:    make(map[string]*model.Dependency),
		retryPolicies:   make(map[string]*model.RetryPolicy),
		circuitBreakers: make(map[string]*model.CircuitBreaker),
		parallelConfigs: make(map[string]*model.ParallelConfig),
		conditions:      make(map[string]*model.Condition),
		auditLogs:       make(map[string]*model.AuditLog),
		envVars:         make(map[string]*model.EnvVar),
		webhooks:        make(map[string]*model.Webhook),
		templates:       make(map[string]*model.Template),
		webhookLogs:     make(map[string]*model.WebhookLog),
	}
}

var _ Store = (*MemoryStore)(nil)
