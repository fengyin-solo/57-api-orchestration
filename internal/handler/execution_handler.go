package handler

import (
	"encoding/json"
	"net/http"

	"orchestr/internal/model"
	"orchestr/pkg/httpx"
)

func (s *Server) registerExecutionRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/executions", s.createExecution)
	mux.HandleFunc("GET /api/executions", s.listExecutions)
	mux.HandleFunc("GET /api/executions/{id}", s.getExecution)
	mux.HandleFunc("PUT /api/executions/{id}", s.updateExecution)
	mux.HandleFunc("DELETE /api/executions/{id}", s.deleteExecution)
	mux.HandleFunc("POST /api/executions/{id}/run", s.runExecution)
	mux.HandleFunc("POST /api/executions/{id}/complete", s.completeExecution)
	mux.HandleFunc("POST /api/executions/{id}/fail", s.failExecution)
	mux.HandleFunc("POST /api/executions/{id}/timeout", s.timeoutExecution)
	mux.HandleFunc("POST /api/executions/execute-flow", s.executeFlow)
	mux.HandleFunc("POST /api/executions/batch-execute", s.batchExecute)
}

type createExecutionRequest struct {
	FlowID string          `json:"flow_id"`
	Input  json.RawMessage `json:"input"`
}

func (s *Server) createExecution(w http.ResponseWriter, r *http.Request) {
	var req createExecutionRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	e, err := s.svc.CreateExecution(model.Execution{FlowID: req.FlowID, Input: req.Input})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, e)
}

func (s *Server) listExecutions(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ExecutionFilter{
		FlowID: r.URL.Query().Get("flow_id"),
		Status: r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListExecutions(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getExecution(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	e, err := s.svc.GetExecution(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, e)
}

func (s *Server) updateExecution(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req createExecutionRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	e, err := s.svc.UpdateExecution(id, model.Execution{FlowID: req.FlowID, Input: req.Input})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, e)
}

func (s *Server) deleteExecution(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteExecution(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) runExecution(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	e, err := s.svc.RunExecution(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, e)
}

func (s *Server) completeExecution(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req struct {
		Output json.RawMessage `json:"output"`
	}
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	e, err := s.svc.CompleteExecution(id, req.Output)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, e)
}

func (s *Server) failExecution(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req struct {
		ErrorMsg string `json:"error_msg"`
	}
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	e, err := s.svc.FailExecution(id, req.ErrorMsg)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, e)
}

func (s *Server) timeoutExecution(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	e, err := s.svc.TimeoutExecution(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, e)
}

type executeFlowRequest struct {
	FlowID string          `json:"flow_id"`
	Input  json.RawMessage `json:"input"`
}

func (s *Server) executeFlow(w http.ResponseWriter, r *http.Request) {
	var req executeFlowRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	e, err := s.svc.ExecuteFlow(req.FlowID, req.Input)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, e)
}

type batchExecuteRequest struct {
	FlowID string            `json:"flow_id"`
	Inputs []json.RawMessage `json:"inputs"`
}

func (s *Server) batchExecute(w http.ResponseWriter, r *http.Request) {
	var req batchExecuteRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	items, err := s.svc.BatchExecute(req.FlowID, req.Inputs)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, items)
}
