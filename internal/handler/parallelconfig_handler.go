package handler

import (
	"net/http"

	"orchestr/internal/model"
	"orchestr/pkg/httpx"
)

func (s *Server) registerParallelConfigRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/parallel-configs", s.createParallelConfig)
	mux.HandleFunc("GET /api/parallel-configs", s.listParallelConfigs)
	mux.HandleFunc("GET /api/parallel-configs/{id}", s.getParallelConfig)
	mux.HandleFunc("PUT /api/parallel-configs/{id}", s.updateParallelConfig)
	mux.HandleFunc("DELETE /api/parallel-configs/{id}", s.deleteParallelConfig)
}

type createParallelConfigRequest struct {
	FlowID  string   `json:"flow_id"`
	StepIDs []string `json:"step_ids"`
	Status  string   `json:"status"`
}

func (s *Server) createParallelConfig(w http.ResponseWriter, r *http.Request) {
	var req createParallelConfigRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	p, err := s.svc.CreateParallelConfig(model.ParallelConfig{FlowID: req.FlowID, StepIDs: req.StepIDs, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, p)
}

func (s *Server) listParallelConfigs(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ParallelConfigFilter{
		FlowID: r.URL.Query().Get("flow_id"),
		Status: r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListParallelConfigs(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getParallelConfig(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	p, err := s.svc.GetParallelConfig(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, p)
}

func (s *Server) updateParallelConfig(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req createParallelConfigRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	p, err := s.svc.UpdateParallelConfig(id, model.ParallelConfig{FlowID: req.FlowID, StepIDs: req.StepIDs, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, p)
}

func (s *Server) deleteParallelConfig(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteParallelConfig(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
