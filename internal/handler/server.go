// Package handler 实现 HTTP 处理器层。
package handler

import (
	"errors"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	"orchestr/internal/config"
	"orchestr/internal/model"
	"orchestr/internal/service"
	"orchestr/internal/store"
	"orchestr/pkg/httpx"
	"orchestr/pkg/logger"
)

type Server struct {
	svc *service.Service
	log *logger.Logger
	cfg *config.Config
}

func NewServer(svc *service.Service, log *logger.Logger, cfg *config.Config) *Server {
	return &Server{svc: svc, log: log, cfg: cfg}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	s.registerServiceRoutes(mux)
	s.registerFlowRoutes(mux)
	s.registerStepRoutes(mux)
	s.registerExecutionRoutes(mux)
	s.registerStepExecutionRoutes(mux)
	s.registerDependencyRoutes(mux)
	s.registerRetryPolicyRoutes(mux)
	s.registerCircuitBreakerRoutes(mux)
	s.registerParallelConfigRoutes(mux)
	s.registerConditionRoutes(mux)
	s.registerAuditLogRoutes(mux)
	s.registerEnvVarRoutes(mux)
	s.registerWebhookRoutes(mux)
	s.registerWebhookLogRoutes(mux)
	s.registerTemplateRoutes(mux)
	s.registerStatsRoutes(mux)
	s.registerExportRoutes(mux)
	s.registerImportRoutes(mux)
	mux.Handle("GET /", http.FileServer(http.Dir("web")))
	return s.corsMiddleware(s.apiKeyMiddleware(s.rateLimitMiddleware(s.loggingMiddleware(s.recoveryMiddleware(mux)))))
}

func (s *Server) maxPageSize() int {
	if s.cfg != nil && s.cfg.MaxPageSize > 0 {
		return s.cfg.MaxPageSize
	}
	return 100
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		s.log.Infof("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func (s *Server) recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				s.log.Errorf("panic: %v\n%s", rec, debug.Stack())
				httpx.InternalError(w, "服务器内部错误")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case model.IsValidationError(err):
		httpx.BadRequest(w, err.Error())
	case errors.Is(err, store.ErrNotFound):
		httpx.NotFound(w, err.Error())
	case errors.Is(err, store.ErrConflict):
		httpx.Conflict(w, err.Error())
	default:
		httpx.InternalError(w, err.Error())
	}
}

var ipRequests = make(map[string][]time.Time)
var ipMu sync.Mutex

func (s *Server) rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		ipMu.Lock()
		now := time.Now()
		window := now.Add(-time.Minute)
		list := ipRequests[ip]
		valid := make([]time.Time, 0, len(list))
		for _, t := range list {
			if t.After(window) {
				valid = append(valid, t)
			}
		}
		valid = append(valid, now)
		ipRequests[ip] = valid
		ipMu.Unlock()
		if len(valid) > 300 {
			httpx.Error(w, http.StatusTooManyRequests, 429, "请求过于频繁")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) apiKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" || r.URL.Path == "/index.html" || r.URL.Path == "/app.js" || r.URL.Path == "/style.css" {
			next.ServeHTTP(w, r)
			return
		}
		key := r.Header.Get("X-API-Key")
		if key == "" {
			key = r.URL.Query().Get("api_key")
		}
		if key != s.cfg.APIKey {
			httpx.Unauthorized(w, "API Key 无效")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-API-Key")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
