package handler

import (
	"encoding/json"
	"net/http"

	"orchestr/internal/model"
	"orchestr/pkg/httpx"
)

func (s *Server) registerStepExecutionRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/step-executions", s.createStepExecution)
	mux.HandleFunc("GET /api/step-executions", s.listStepExecutions)
	mux.HandleFunc("GET /api/step-executions/{id}", s.getStepExecution)
	mux.HandleFunc("PUT /api/step-executions/{id}", s.updateStepExecution)
	mux.HandleFunc("DELETE /api/step-executions/{id}", s.deleteStepExecution)
	mux.HandleFunc("POST /api/step-executions/{id}/transition", s.transitionStepExecution)
}

type createStepExecutionRequest struct {
	ExecutionID string          `json:"execution_id"`
	StepID      string          `json:"step_id"`
	Status      string          `json:"status"`
	Response    json.RawMessage `json:"response"`
	DurationMs  int             `json:"duration_ms"`
}

func (s *Server) createStepExecution(w http.ResponseWriter, r *http.Request) {
	var req createStepExecutionRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	se, err := s.svc.CreateStepExecution(model.StepExecution{ExecutionID: req.ExecutionID, StepID: req.StepID, Status: req.Status, Response: req.Response, DurationMs: req.DurationMs})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, se)
}

func (s *Server) listStepExecutions(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.StepExecutionFilter{
		ExecutionID: r.URL.Query().Get("execution_id"),
		Status:      r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListStepExecutions(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getStepExecution(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	se, err := s.svc.GetStepExecution(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, se)
}

func (s *Server) updateStepExecution(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req createStepExecutionRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	se, err := s.svc.UpdateStepExecution(id, model.StepExecution{Status: req.Status, Response: req.Response, DurationMs: req.DurationMs})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, se)
}

func (s *Server) deleteStepExecution(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteStepExecution(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) transitionStepExecution(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req struct {
		Status string `json:"status"`
	}
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	se, err := s.svc.TransitionStepExecution(id, req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, se)
}
