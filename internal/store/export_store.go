package store

func (s *MemoryStore) ExportAll() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make(map[string]interface{})
	result["services"] = s.services
	result["flows"] = s.flows
	result["steps"] = s.steps
	result["executions"] = s.executions
	result["stepExecutions"] = s.stepExecutions
	result["dependencies"] = s.dependencies
	result["retryPolicies"] = s.retryPolicies
	result["circuitBreakers"] = s.circuitBreakers
	result["parallelConfigs"] = s.parallelConfigs
	result["conditions"] = s.conditions
	result["auditLogs"] = s.auditLogs
	result["envVars"] = s.envVars
	result["webhooks"] = s.webhooks
	result["templates"] = s.templates
	result["webhookLogs"] = s.webhookLogs
	return result
}
