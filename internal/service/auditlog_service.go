package service

import (
	"sort"
	"time"

	"orchestr/internal/model"
	"orchestr/pkg/idgen"
)

func (s *Service) logAudit(operator, action, targetType, targetID, detail string) error {
	a := &model.AuditLog{
		ID:         idgen.Hex(),
		Operator:   operator,
		Action:     action,
		TargetType: targetType,
		TargetID:   targetID,
		Detail:     detail,
		CreatedAt:  time.Now().Format(time.RFC3339),
	}
	return s.store.CreateAuditLog(a)
}

func (s *Service) CreateAuditLog(input model.AuditLog) (*model.AuditLog, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	input.ID = idgen.Hex()
	input.CreatedAt = time.Now().Format(time.RFC3339)
	if err := s.store.CreateAuditLog(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *Service) GetAuditLog(id string) (*model.AuditLog, error) {
	return s.store.GetAuditLog(id)
}

func (s *Service) ListAuditLogs(filter model.AuditLogFilter, page, size int) ([]*model.AuditLog, int, error) {
	all := s.store.ListAuditLogs()
	matched := make([]*model.AuditLog, 0, len(all))
	for _, a := range all {
		if filter.Match(a) {
			matched = append(matched, a)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt > matched[j].CreatedAt
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.AuditLog{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) DeleteAuditLog(id string) error {
	return s.store.DeleteAuditLog(id)
}
