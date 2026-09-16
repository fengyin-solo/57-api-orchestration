package handler

import (
	"net/http"

	"orchestr/internal/model"
	"orchestr/pkg/httpx"
)

func (s *Server) registerTemplateRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/templates", s.createTemplate)
	mux.HandleFunc("GET /api/templates", s.listTemplates)
	mux.HandleFunc("GET /api/templates/{id}", s.getTemplate)
	mux.HandleFunc("PUT /api/templates/{id}", s.updateTemplate)
	mux.HandleFunc("DELETE /api/templates/{id}", s.deleteTemplate)
	mux.HandleFunc("POST /api/templates/{id}/render", s.renderTemplate)
}

type createTemplateRequest struct {
	Name        string `json:"name"`
	FlowID      string `json:"flow_id"`
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
	Status      string `json:"status"`
}

func (s *Server) createTemplate(w http.ResponseWriter, r *http.Request) {
	var req createTemplateRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.CreateTemplate(model.Template{Name: req.Name, FlowID: req.FlowID, Content: req.Content, ContentType: req.ContentType, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, t)
}

func (s *Server) listTemplates(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.TemplateFilter{
		FlowID: r.URL.Query().Get("flow_id"),
		Status: r.URL.Query().Get("status"),
		Name:   r.URL.Query().Get("name"),
	}
	items, total, err := s.svc.ListTemplates(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getTemplate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	t, err := s.svc.GetTemplate(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

func (s *Server) updateTemplate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req createTemplateRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.UpdateTemplate(id, model.Template{Name: req.Name, FlowID: req.FlowID, Content: req.Content, ContentType: req.ContentType, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

func (s *Server) deleteTemplate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteTemplate(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) renderTemplate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req struct {
		Vars map[string]string `json:"vars"`
	}
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.GetTemplate(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	result := t.ReplaceVars(req.Vars)
	httpx.OK(w, map[string]string{"result": result})
}
