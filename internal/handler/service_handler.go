package handler

import (
	"net/http"

	"orchestr/internal/model"
	"orchestr/pkg/httpx"
)

func (s *Server) registerServiceRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/services", s.createService)
	mux.HandleFunc("GET /api/services", s.listServices)
	mux.HandleFunc("GET /api/services/{id}", s.getService)
	mux.HandleFunc("PUT /api/services/{id}", s.updateService)
	mux.HandleFunc("DELETE /api/services/{id}", s.deleteService)
}

type createServiceRequest struct {
	Name       string `json:"name"`
	BaseURL    string `json:"base_url"`
	AuthType   string `json:"auth_type"`
	TimeoutMs  int    `json:"timeout_ms"`
	RetryCount int    `json:"retry_count"`
	Status     string `json:"status"`
}

func (s *Server) createService(w http.ResponseWriter, r *http.Request) {
	var req createServiceRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sv, err := s.svc.CreateService(model.Service{Name: req.Name, BaseURL: req.BaseURL, AuthType: req.AuthType, TimeoutMs: req.TimeoutMs, RetryCount: req.RetryCount, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, sv)
}

func (s *Server) listServices(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ServiceFilter{
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListServices(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getService(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	sv, err := s.svc.GetService(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sv)
}

func (s *Server) updateService(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req createServiceRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sv, err := s.svc.UpdateService(id, model.Service{Name: req.Name, BaseURL: req.BaseURL, AuthType: req.AuthType, TimeoutMs: req.TimeoutMs, RetryCount: req.RetryCount, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sv)
}

func (s *Server) deleteService(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteService(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
