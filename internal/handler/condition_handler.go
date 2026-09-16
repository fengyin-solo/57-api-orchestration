package handler

import (
	"net/http"

	"orchestr/internal/model"
	"orchestr/pkg/httpx"
)

func (s *Server) registerConditionRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/conditions", s.createCondition)
	mux.HandleFunc("GET /api/conditions", s.listConditions)
	mux.HandleFunc("GET /api/conditions/{id}", s.getCondition)
	mux.HandleFunc("PUT /api/conditions/{id}", s.updateCondition)
	mux.HandleFunc("DELETE /api/conditions/{id}", s.deleteCondition)
}

type createConditionRequest struct {
	FlowID      string `json:"flow_id"`
	StepID      string `json:"step_id"`
	Expression  string `json:"expression"`
	TrueStepID  string `json:"true_step_id"`
	FalseStepID string `json:"false_step_id"`
	Status      string `json:"status"`
}

func (s *Server) createCondition(w http.ResponseWriter, r *http.Request) {
	var req createConditionRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.CreateCondition(model.Condition{FlowID: req.FlowID, StepID: req.StepID, Expression: req.Expression, TrueStepID: req.TrueStepID, FalseStepID: req.FalseStepID, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, c)
}

func (s *Server) listConditions(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ConditionFilter{
		FlowID: r.URL.Query().Get("flow_id"),
		Status: r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListConditions(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getCondition(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	c, err := s.svc.GetCondition(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

func (s *Server) updateCondition(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req createConditionRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.UpdateCondition(id, model.Condition{FlowID: req.FlowID, StepID: req.StepID, Expression: req.Expression, TrueStepID: req.TrueStepID, FalseStepID: req.FalseStepID, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

func (s *Server) deleteCondition(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteCondition(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
