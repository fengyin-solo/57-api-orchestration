package handler

import (
	"net/http"

	"orchestr/pkg/httpx"
)

func (s *Server) registerStatsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/stats/overview", s.statsOverview)
	mux.HandleFunc("GET /api/stats/success-rate", s.statsSuccessRate)
	mux.HandleFunc("GET /api/stats/by-flow", s.statsByFlow)
	mux.HandleFunc("GET /api/stats/by-service", s.statsByService)
	mux.HandleFunc("GET /api/stats/failure-reasons", s.statsFailureReasons)
}

func (s *Server) statsOverview(w http.ResponseWriter, r *http.Request) {
	overview, err := s.svc.StatsOverview()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, overview)
}

func (s *Server) statsSuccessRate(w http.ResponseWriter, r *http.Request) {
	result := s.svc.StatsSuccessRate()
	httpx.OK(w, result)
}

func (s *Server) statsByFlow(w http.ResponseWriter, r *http.Request) {
	result := s.svc.StatsByFlow()
	httpx.OK(w, result)
}

func (s *Server) statsByService(w http.ResponseWriter, r *http.Request) {
	result := s.svc.StatsByService()
	httpx.OK(w, result)
}

func (s *Server) statsFailureReasons(w http.ResponseWriter, r *http.Request) {
	result := s.svc.StatsFailureReasons()
	httpx.OK(w, result)
}
