package handler

import (
	"net/http"

	"orchestr/internal/model"
	"orchestr/pkg/httpx"
)

func (s *Server) registerEnvVarRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/env-vars", s.createEnvVar)
	mux.HandleFunc("GET /api/env-vars", s.listEnvVars)
	mux.HandleFunc("GET /api/env-vars/{id}", s.getEnvVar)
	mux.HandleFunc("PUT /api/env-vars/{id}", s.updateEnvVar)
	mux.HandleFunc("DELETE /api/env-vars/{id}", s.deleteEnvVar)
	mux.HandleFunc("GET /api/env-vars/flow/{flow_id}", s.getEnvVarsByFlow)
}

type createEnvVarRequest struct {
	FlowID   string `json:"flow_id"`
	Key      string `json:"key"`
	Value    string `json:"value"`
	IsSecret bool   `json:"is_secret"`
	Status   string `json:"status"`
}

func (s *Server) createEnvVar(w http.ResponseWriter, r *http.Request) {
	var req createEnvVarRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	e, err := s.svc.CreateEnvVar(model.EnvVar{FlowID: req.FlowID, Key: req.Key, Value: req.Value, IsSecret: req.IsSecret, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, e)
}

func (s *Server) listEnvVars(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.EnvVarFilter{
		FlowID: r.URL.Query().Get("flow_id"),
		Status: r.URL.Query().Get("status"),
		Key:    r.URL.Query().Get("key"),
	}
	items, total, err := s.svc.ListEnvVars(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getEnvVar(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	e, err := s.svc.GetEnvVar(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, e)
}

func (s *Server) updateEnvVar(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req createEnvVarRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	e, err := s.svc.UpdateEnvVar(id, model.EnvVar{FlowID: req.FlowID, Key: req.Key, Value: req.Value, IsSecret: req.IsSecret, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, e)
}

func (s *Server) deleteEnvVar(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteEnvVar(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) getEnvVarsByFlow(w http.ResponseWriter, r *http.Request) {
	flowID := r.PathValue("flow_id")
	m := s.svc.BuildEnvMap(flowID)
	httpx.OK(w, m)
}
