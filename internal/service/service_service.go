package service

import (
	"sort"
	"time"

	"orchestr/internal/model"
	"orchestr/pkg/idgen"
)

func (s *Service) CreateService(input model.Service) (*model.Service, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now().Format(time.RFC3339)
	sv := &model.Service{
		ID:         idgen.Hex(),
		Name:       input.Name,
		BaseURL:    input.BaseURL,
		AuthType:   input.AuthType,
		TimeoutMs:  input.TimeoutMs,
		RetryCount: input.RetryCount,
		Status:     input.Status,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := s.store.CreateService(sv); err != nil {
		return nil, err
	}
	_ = s.logAudit("system", "create_service", "service", sv.ID, "创建服务 "+sv.Name)
	return sv, nil
}

func (s *Service) GetService(id string) (*model.Service, error) {
	return s.store.GetService(id)
}

func (s *Service) ListServices(filter model.ServiceFilter, page, size int) ([]*model.Service, int, error) {
	all := s.store.ListServices()
	matched := make([]*model.Service, 0, len(all))
	for _, sv := range all {
		if filter.Match(sv) {
			matched = append(matched, sv)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt > matched[j].CreatedAt
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Service{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateService(id string, input model.Service) (*model.Service, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	sv, err := s.store.GetService(id)
	if err != nil {
		return nil, err
	}
	sv.Name = input.Name
	sv.BaseURL = input.BaseURL
	sv.AuthType = input.AuthType
	sv.TimeoutMs = input.TimeoutMs
	sv.RetryCount = input.RetryCount
	sv.Status = input.Status
	sv.UpdatedAt = time.Now().Format(time.RFC3339)
	if err := s.store.UpdateService(sv); err != nil {
		return nil, err
	}
	_ = s.logAudit("system", "update_service", "service", sv.ID, "更新服务 "+sv.Name)
	return sv, nil
}

func (s *Service) DeleteService(id string) error {
	if err := s.store.DeleteService(id); err != nil {
		return err
	}
	_ = s.logAudit("system", "delete_service", "service", id, "删除服务")
	return nil
}
