package handler

import (
	"net/http"

	"orchestr/internal/model"
	"orchestr/pkg/httpx"
)

func (s *Server) registerCircuitBreakerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/circuit-breakers", s.createCircuitBreaker)
	mux.HandleFunc("GET /api/circuit-breakers", s.listCircuitBreakers)
	mux.HandleFunc("GET /api/circuit-breakers/{id}", s.getCircuitBreaker)
	mux.HandleFunc("PUT /api/circuit-breakers/{id}", s.updateCircuitBreaker)
	mux.HandleFunc("DELETE /api/circuit-breakers/{id}", s.deleteCircuitBreaker)
}

type createCircuitBreakerRequest struct {
	ServiceID  string `json:"service_id"`
	Threshold  int    `json:"threshold"`
	CooldownMs int    `json:"cooldown_ms"`
	Status     string `json:"status"`
}

func (s *Server) createCircuitBreaker(w http.ResponseWriter, r *http.Request) {
	var req createCircuitBreakerRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.CreateCircuitBreaker(model.CircuitBreaker{ServiceID: req.ServiceID, Threshold: req.Threshold, CooldownMs: req.CooldownMs, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, c)
}

func (s *Server) listCircuitBreakers(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.CircuitBreakerFilter{
		ServiceID: r.URL.Query().Get("service_id"),
		Status:    r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListCircuitBreakers(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getCircuitBreaker(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	c, err := s.svc.GetCircuitBreaker(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

func (s *Server) updateCircuitBreaker(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req createCircuitBreakerRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.UpdateCircuitBreaker(id, model.CircuitBreaker{ServiceID: req.ServiceID, Threshold: req.Threshold, CooldownMs: req.CooldownMs, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

func (s *Server) deleteCircuitBreaker(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteCircuitBreaker(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
