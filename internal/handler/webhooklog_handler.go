package handler

import (
	"net/http"

	"orchestr/internal/model"
	"orchestr/pkg/httpx"
)

func (s *Server) registerWebhookLogRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/webhook-logs", s.createWebhookLog)
	mux.HandleFunc("GET /api/webhook-logs", s.listWebhookLogs)
	mux.HandleFunc("GET /api/webhook-logs/{id}", s.getWebhookLog)
	mux.HandleFunc("PUT /api/webhook-logs/{id}", s.updateWebhookLog)
	mux.HandleFunc("DELETE /api/webhook-logs/{id}", s.deleteWebhookLog)
}

type createWebhookLogRequest struct {
	WebhookID string `json:"webhook_id"`
	Payload   string `json:"payload"`
	Status    string `json:"status"`
	Response  string `json:"response"`
}

func (s *Server) createWebhookLog(w http.ResponseWriter, r *http.Request) {
	var req createWebhookLogRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	wl, err := s.svc.CreateWebhookLog(model.WebhookLog{WebhookID: req.WebhookID, Payload: req.Payload, Status: req.Status, Response: req.Response})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, wl)
}

func (s *Server) listWebhookLogs(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.WebhookLogFilter{
		WebhookID: r.URL.Query().Get("webhook_id"),
		Status:    r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListWebhookLogs(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getWebhookLog(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	wl, err := s.svc.GetWebhookLog(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, wl)
}

func (s *Server) updateWebhookLog(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req createWebhookLogRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	wl, err := s.svc.UpdateWebhookLog(id, model.WebhookLog{WebhookID: req.WebhookID, Payload: req.Payload, Status: req.Status, Response: req.Response})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, wl)
}

func (s *Server) deleteWebhookLog(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteWebhookLog(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
