package handler

import (
	"encoding/json"
	"net/http"

	"orchestr/internal/model"
	"orchestr/pkg/httpx"
)

func (s *Server) registerStepRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/steps", s.createStep)
	mux.HandleFunc("GET /api/steps", s.listSteps)
	mux.HandleFunc("GET /api/steps/{id}", s.getStep)
	mux.HandleFunc("PUT /api/steps/{id}", s.updateStep)
	mux.HandleFunc("DELETE /api/steps/{id}", s.deleteStep)
	mux.HandleFunc("POST /api/steps/batch", s.batchCreateSteps)
}

type createStepRequest struct {
	FlowID        string          `json:"flow_id"`
	ServiceID     string          `json:"service_id"`
	Method        string          `json:"method"`
	Path          string          `json:"path"`
	ParamsMapping json.RawMessage `json:"params_mapping"`
	Order         int             `json:"order"`
	TimeoutMs     int             `json:"timeout_ms"`
	RetryCount    int             `json:"retry_count"`
	Status        string          `json:"status"`
}

func (s *Server) createStep(w http.ResponseWriter, r *http.Request) {
	var req createStepRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	st, err := s.svc.CreateStep(model.Step{FlowID: req.FlowID, ServiceID: req.ServiceID, Method: req.Method, Path: req.Path, ParamsMapping: req.ParamsMapping, Order: req.Order, TimeoutMs: req.TimeoutMs, RetryCount: req.RetryCount, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, st)
}

func (s *Server) listSteps(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.StepFilter{
		FlowID:    r.URL.Query().Get("flow_id"),
		Status:    r.URL.Query().Get("status"),
		Method:    r.URL.Query().Get("method"),
		ServiceID: r.URL.Query().Get("service_id"),
	}
	items, total, err := s.svc.ListSteps(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getStep(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	st, err := s.svc.GetStep(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, st)
}

func (s *Server) updateStep(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req createStepRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	st, err := s.svc.UpdateStep(id, model.Step{FlowID: req.FlowID, ServiceID: req.ServiceID, Method: req.Method, Path: req.Path, ParamsMapping: req.ParamsMapping, Order: req.Order, TimeoutMs: req.TimeoutMs, RetryCount: req.RetryCount, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, st)
}

func (s *Server) deleteStep(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteStep(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type batchCreateStepsRequest struct {
	Steps []createStepRequest `json:"steps"`
}

func (s *Server) batchCreateSteps(w http.ResponseWriter, r *http.Request) {
	var req batchCreateStepsRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	inputs := make([]model.Step, 0, len(req.Steps))
	for _, rs := range req.Steps {
		inputs = append(inputs, model.Step{FlowID: rs.FlowID, ServiceID: rs.ServiceID, Method: rs.Method, Path: rs.Path, ParamsMapping: rs.ParamsMapping, Order: rs.Order, TimeoutMs: rs.TimeoutMs, RetryCount: rs.RetryCount, Status: rs.Status})
	}
	items, err := s.svc.BatchCreateSteps(inputs)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, items)
}
