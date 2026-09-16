package handler

import (
	"net/http"

	"orchestr/pkg/httpx"
)

func (s *Server) registerExportRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/export", s.exportAll)
}

func (s *Server) exportAll(w http.ResponseWriter, r *http.Request) {
	data := s.svc.ExportAll()
	httpx.OK(w, data)
}
