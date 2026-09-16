package service

import (
	"sort"
	"time"

	"orchestr/internal/model"
	"orchestr/pkg/idgen"
)

func (s *Service) CreateTemplate(input model.Template) (*model.Template, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now().Format(time.RFC3339)
	t := &model.Template{
		ID:          idgen.Hex(),
		Name:        input.Name,
		FlowID:      input.FlowID,
		Content:     input.Content,
		ContentType: input.ContentType,
		Status:      input.Status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.store.CreateTemplate(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) GetTemplate(id string) (*model.Template, error) {
	return s.store.GetTemplate(id)
}

func (s *Service) ListTemplates(filter model.TemplateFilter, page, size int) ([]*model.Template, int, error) {
	all := s.store.ListTemplates()
	matched := make([]*model.Template, 0, len(all))
	for _, t := range all {
		if filter.Match(t) {
			matched = append(matched, t)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt > matched[j].CreatedAt
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Template{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateTemplate(id string, input model.Template) (*model.Template, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	t, err := s.store.GetTemplate(id)
	if err != nil {
		return nil, err
	}
	t.Name = input.Name
	t.FlowID = input.FlowID
	t.Content = input.Content
	t.ContentType = input.ContentType
	t.Status = input.Status
	t.UpdatedAt = time.Now().Format(time.RFC3339)
	if err := s.store.UpdateTemplate(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) DeleteTemplate(id string) error {
	return s.store.DeleteTemplate(id)
}

func (s *Service) RenderTemplate(name string, vars map[string]string) (string, error) {
	t, err := s.store.GetTemplateByName(name)
	if err != nil {
		return "", err
	}
	return t.ReplaceVars(vars), nil
}
