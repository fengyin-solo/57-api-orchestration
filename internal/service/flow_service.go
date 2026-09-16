package service

import (
	"sort"
	"time"

	"orchestr/internal/model"
	"orchestr/pkg/idgen"
)

func (s *Service) CreateFlow(input model.Flow) (*model.Flow, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now().Format(time.RFC3339)
	f := &model.Flow{
		ID:          idgen.Hex(),
		Name:        input.Name,
		Description: input.Description,
		Status:      input.Status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.store.CreateFlow(f); err != nil {
		return nil, err
	}
	_ = s.logAudit("system", "create_flow", "flow", f.ID, "创建流程 "+f.Name)
	return f, nil
}

func (s *Service) GetFlow(id string) (*model.Flow, error) {
	return s.store.GetFlow(id)
}

func (s *Service) ListFlows(filter model.FlowFilter, page, size int) ([]*model.Flow, int, error) {
	all := s.store.ListFlows()
	matched := make([]*model.Flow, 0, len(all))
	for _, f := range all {
		if filter.Match(f) {
			matched = append(matched, f)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt > matched[j].CreatedAt
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Flow{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateFlow(id string, input model.Flow) (*model.Flow, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	f, err := s.store.GetFlow(id)
	if err != nil {
		return nil, err
	}
	f.Name = input.Name
	f.Description = input.Description
	f.Status = input.Status
	f.UpdatedAt = time.Now().Format(time.RFC3339)
	if err := s.store.UpdateFlow(f); err != nil {
		return nil, err
	}
	_ = s.logAudit("system", "update_flow", "flow", f.ID, "更新流程 "+f.Name)
	return f, nil
}

func (s *Service) DeleteFlow(id string) error {
	if err := s.store.DeleteFlow(id); err != nil {
		return err
	}
	_ = s.logAudit("system", "delete_flow", "flow", id, "删除流程")
	return nil
}
