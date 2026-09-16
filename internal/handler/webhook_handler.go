package handler

import (
	"net/http"

	"orchestr/internal/model"
	"orchestr/pkg/httpx"
)

func (s *Server) registerWebhookRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/webhooks", s.createWebhook)
	mux.HandleFunc("GET /api/webhooks", s.listWebhooks)
	mux.HandleFunc("GET /api/webhooks/{id}", s.getWebhook)
	mux.HandleFunc("PUT /api/webhooks/{id}", s.updateWebhook)
	mux.HandleFunc("DELETE /api/webhooks/{id}", s.deleteWebhook)
	mux.HandleFunc("POST /api/webhooks/{id}/trigger", s.triggerWebhook)
}

type createWebhookRequest struct {
	FlowID     string `json:"flow_id"`
	URL        string `json:"url"`
	Method     string `json:"method"`
	Headers    string `json:"headers"`
	RetryCount int    `json:"retry_count"`
	TimeoutMs  int    `json:"timeout_ms"`
	Status     string `json:"status"`
}

func (s *Server) createWebhook(w http.ResponseWriter, r *http.Request) {
	var req createWebhookRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	wh, err := s.svc.CreateWebhook(model.Webhook{FlowID: req.FlowID, URL: req.URL, Method: req.Method, Headers: req.Headers, RetryCount: req.RetryCount, TimeoutMs: req.TimeoutMs, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, wh)
}

func (s *Server) listWebhooks(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.WebhookFilter{
		FlowID: r.URL.Query().Get("flow_id"),
		Status: r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListWebhooks(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getWebhook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	wh, err := s.svc.GetWebhook(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, wh)
}

func (s *Server) updateWebhook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req createWebhookRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	wh, err := s.svc.UpdateWebhook(id, model.Webhook{FlowID: req.FlowID, URL: req.URL, Method: req.Method, Headers: req.Headers, RetryCount: req.RetryCount, TimeoutMs: req.TimeoutMs, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, wh)
}

func (s *Server) deleteWebhook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteWebhook(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) triggerWebhook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	wh, err := s.svc.GetWebhook(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	if err := s.svc.TriggerWebhooks(wh.FlowID, "{}"); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]string{"status": "triggered"})
}
