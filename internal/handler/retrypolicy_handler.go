package handler

import (
	"net/http"

	"orchestr/internal/model"
	"orchestr/pkg/httpx"
)

func (s *Server) registerRetryPolicyRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/retry-policies", s.createRetryPolicy)
	mux.HandleFunc("GET /api/retry-policies", s.listRetryPolicies)
	mux.HandleFunc("GET /api/retry-policies/{id}", s.getRetryPolicy)
	mux.HandleFunc("PUT /api/retry-policies/{id}", s.updateRetryPolicy)
	mux.HandleFunc("DELETE /api/retry-policies/{id}", s.deleteRetryPolicy)
}

type createRetryPolicyRequest struct {
	Name           string `json:"name"`
	MaxRetries     int    `json:"max_retries"`
	BackoffType    string `json:"backoff_type"`
	InitialDelayMs int    `json:"initial_delay_ms"`
	Status         string `json:"status"`
}

func (s *Server) createRetryPolicy(w http.ResponseWriter, r *http.Request) {
	var req createRetryPolicyRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rp, err := s.svc.CreateRetryPolicy(model.RetryPolicy{Name: req.Name, MaxRetries: req.MaxRetries, BackoffType: req.BackoffType, InitialDelayMs: req.InitialDelayMs, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, rp)
}

func (s *Server) listRetryPolicies(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.RetryPolicyFilter{
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListRetryPolicies(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getRetryPolicy(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	rp, err := s.svc.GetRetryPolicy(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rp)
}

func (s *Server) updateRetryPolicy(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req createRetryPolicyRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rp, err := s.svc.UpdateRetryPolicy(id, model.RetryPolicy{Name: req.Name, MaxRetries: req.MaxRetries, BackoffType: req.BackoffType, InitialDelayMs: req.InitialDelayMs, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rp)
}

func (s *Server) deleteRetryPolicy(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteRetryPolicy(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
