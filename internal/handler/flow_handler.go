package handler

import (
	"net/http"

	"orchestr/internal/model"
	"orchestr/pkg/httpx"
)

func (s *Server) registerFlowRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/flows", s.createFlow)
	mux.HandleFunc("GET /api/flows", s.listFlows)
	mux.HandleFunc("GET /api/flows/{id}", s.getFlow)
	mux.HandleFunc("PUT /api/flows/{id}", s.updateFlow)
	mux.HandleFunc("DELETE /api/flows/{id}", s.deleteFlow)
}

type createFlowRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func (s *Server) createFlow(w http.ResponseWriter, r *http.Request) {
	var req createFlowRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	f, err := s.svc.CreateFlow(model.Flow{Name: req.Name, Description: req.Description, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, f)
}

func (s *Server) listFlows(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.FlowFilter{
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListFlows(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getFlow(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	f, err := s.svc.GetFlow(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, f)
}

func (s *Server) updateFlow(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req createFlowRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	f, err := s.svc.UpdateFlow(id, model.Flow{Name: req.Name, Description: req.Description, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, f)
}

func (s *Server) deleteFlow(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteFlow(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
