package handler

import (
	"net/http"

	"orchestr/pkg/httpx"
)

func (s *Server) registerImportRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/import", s.importSnapshot)
	mux.HandleFunc("POST /api/import/validate", s.validateSnapshot)
}

type importSnapshotRequest struct {
	Data map[string]interface{} `json:"data"`
}

func (s *Server) importSnapshot(w http.ResponseWriter, r *http.Request) {
	var req importSnapshotRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	if err := s.svc.ValidateSnapshot(req.Data); err != nil {
		httpx.BadRequest(w, err.Error())
		return
	}
	if err := s.svc.ImportSnapshot(req.Data); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, map[string]string{"status": "imported"})
}

func (s *Server) validateSnapshot(w http.ResponseWriter, r *http.Request) {
	var req importSnapshotRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	if err := s.svc.ValidateSnapshot(req.Data); err != nil {
		httpx.BadRequest(w, err.Error())
		return
	}
	httpx.OK(w, map[string]string{"status": "valid"})
}
