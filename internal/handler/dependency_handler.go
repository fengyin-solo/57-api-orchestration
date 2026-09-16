package handler

import (
	"net/http"

	"orchestr/internal/model"
	"orchestr/pkg/httpx"
)

func (s *Server) registerDependencyRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/dependencies", s.createDependency)
	mux.HandleFunc("GET /api/dependencies", s.listDependencies)
	mux.HandleFunc("GET /api/dependencies/{id}", s.getDependency)
	mux.HandleFunc("PUT /api/dependencies/{id}", s.updateDependency)
	mux.HandleFunc("DELETE /api/dependencies/{id}", s.deleteDependency)
}

type createDependencyRequest struct {
	FlowID          string `json:"flow_id"`
	StepID          string `json:"step_id"`
	DependsOnStepID string `json:"depends_on_step_id"`
	Condition       string `json:"condition"`
}

func (s *Server) createDependency(w http.ResponseWriter, r *http.Request) {
	var req createDependencyRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	d, err := s.svc.CreateDependency(model.Dependency{FlowID: req.FlowID, StepID: req.StepID, DependsOnStepID: req.DependsOnStepID, Condition: req.Condition})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, d)
}

func (s *Server) listDependencies(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.DependencyFilter{
		FlowID:          r.URL.Query().Get("flow_id"),
		StepID:          r.URL.Query().Get("step_id"),
		DependsOnStepID: r.URL.Query().Get("depends_on_step_id"),
	}
	items, total, err := s.svc.ListDependencies(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getDependency(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	d, err := s.svc.GetDependency(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, d)
}

func (s *Server) updateDependency(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req createDependencyRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	d, err := s.svc.UpdateDependency(id, model.Dependency{FlowID: req.FlowID, StepID: req.StepID, DependsOnStepID: req.DependsOnStepID, Condition: req.Condition})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, d)
}

func (s *Server) deleteDependency(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteDependency(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
