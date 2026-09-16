package service

import (
	"encoding/json"
	"fmt"

	"orchestr/internal/model"
)

func (s *Service) ImportSnapshot(data map[string]interface{}) error {
	if svcs, ok := data["services"].([]interface{}); ok {
		for _, v := range svcs {
			var sv model.Service
			b, _ := json.Marshal(v)
			_ = json.Unmarshal(b, &sv)
			if sv.ID != "" {
				_ = s.store.CreateService(&sv)
			}
		}
	}
	if flows, ok := data["flows"].([]interface{}); ok {
		for _, v := range flows {
			var f model.Flow
			b, _ := json.Marshal(v)
			_ = json.Unmarshal(b, &f)
			if f.ID != "" {
				_ = s.store.CreateFlow(&f)
			}
		}
	}
	if steps, ok := data["steps"].([]interface{}); ok {
		for _, v := range steps {
			var st model.Step
			b, _ := json.Marshal(v)
			_ = json.Unmarshal(b, &st)
			if st.ID != "" {
				_ = s.store.CreateStep(&st)
			}
		}
	}
	if deps, ok := data["dependencies"].([]interface{}); ok {
		for _, v := range deps {
			var d model.Dependency
			b, _ := json.Marshal(v)
			_ = json.Unmarshal(b, &d)
			if d.ID != "" {
				_ = s.store.CreateDependency(&d)
			}
		}
	}
	if cbs, ok := data["circuitBreakers"].([]interface{}); ok {
		for _, v := range cbs {
			var c model.CircuitBreaker
			b, _ := json.Marshal(v)
			_ = json.Unmarshal(b, &c)
			if c.ID != "" {
				_ = s.store.CreateCircuitBreaker(&c)
			}
		}
	}
	if rps, ok := data["retryPolicies"].([]interface{}); ok {
		for _, v := range rps {
			var rp model.RetryPolicy
			b, _ := json.Marshal(v)
			_ = json.Unmarshal(b, &rp)
			if rp.ID != "" {
				_ = s.store.CreateRetryPolicy(&rp)
			}
		}
	}
	if pcs, ok := data["parallelConfigs"].([]interface{}); ok {
		for _, v := range pcs {
			var p model.ParallelConfig
			b, _ := json.Marshal(v)
			_ = json.Unmarshal(b, &p)
			if p.ID != "" {
				_ = s.store.CreateParallelConfig(&p)
			}
		}
	}
	if conds, ok := data["conditions"].([]interface{}); ok {
		for _, v := range conds {
			var c model.Condition
			b, _ := json.Marshal(v)
			_ = json.Unmarshal(b, &c)
			if c.ID != "" {
				_ = s.store.CreateCondition(&c)
			}
		}
	}
	if evs, ok := data["envVars"].([]interface{}); ok {
		for _, v := range evs {
			var e model.EnvVar
			b, _ := json.Marshal(v)
			_ = json.Unmarshal(b, &e)
			if e.ID != "" {
				_ = s.store.CreateEnvVar(&e)
			}
		}
	}
	if whs, ok := data["webhooks"].([]interface{}); ok {
		for _, v := range whs {
			var w model.Webhook
			b, _ := json.Marshal(v)
			_ = json.Unmarshal(b, &w)
			if w.ID != "" {
				_ = s.store.CreateWebhook(&w)
			}
		}
	}
	if tpls, ok := data["templates"].([]interface{}); ok {
		for _, v := range tpls {
			var t model.Template
			b, _ := json.Marshal(v)
			_ = json.Unmarshal(b, &t)
			if t.ID != "" {
				_ = s.store.CreateTemplate(&t)
			}
		}
	}
	_ = s.logAudit("system", "import_snapshot", "system", "all", "导入快照")
	return nil
}

func (s *Service) ValidateSnapshot(data map[string]interface{}) error {
	if data == nil {
		return fmt.Errorf("数据不能为空")
	}
	return nil
}
