package handler

import (
	"net/http"

	"lightweightrpc/internal/model"
	"lightweightrpc/pkg/httpx"
)

func (s *Server) registerRetryPolicyRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/retry-policies", s.createRetryPolicy)
	mux.HandleFunc("GET /api/retry-policies", s.listRetryPolicy)
	mux.HandleFunc("GET /api/retry-policies/{id}", s.getRetryPolicy)
	mux.HandleFunc("PUT /api/retry-policies/{id}", s.updateRetryPolicy)
	mux.HandleFunc("DELETE /api/retry-policies/{id}", s.deleteRetryPolicy)
	mux.HandleFunc("PATCH /api/retry-policies/{id}/status", s.transitionRetryPolicy)
}

type createRetryPolicyRequest struct {
	Name          string  `json:"name"`
	MaxRetries    int     `json:"max_retries"`
	BaseTimeoutMs int     `json:"base_timeout_ms"`
	Backoff       string  `json:"backoff"`
	Jitter        float64 `json:"jitter"`
}

func (s *Server) createRetryPolicy(w http.ResponseWriter, r *http.Request) {
	var req createRetryPolicyRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	x, err := s.svc.CreateRetryPolicy(model.RetryPolicy{Name: req.Name, MaxRetries: req.MaxRetries, BaseTimeoutMs: req.BaseTimeoutMs, Backoff: req.Backoff, Jitter: req.Jitter})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, x)
}

func (s *Server) listRetryPolicy(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.RetryPolicyFilter{
		Backoff: r.URL.Query().Get("backoff"),
		Status:  r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListRetryPolicys(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getRetryPolicy(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	x, err := s.svc.GetRetryPolicy(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, x)
}

type updateRetryPolicyRequest struct {
	Name          string  `json:"name"`
	MaxRetries    int     `json:"max_retries"`
	BaseTimeoutMs int     `json:"base_timeout_ms"`
	Backoff       string  `json:"backoff"`
	Jitter        float64 `json:"jitter"`
}

func (s *Server) updateRetryPolicy(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateRetryPolicyRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	x, err := s.svc.UpdateRetryPolicy(id, model.RetryPolicy{Name: req.Name, MaxRetries: req.MaxRetries, BaseTimeoutMs: req.BaseTimeoutMs, Backoff: req.Backoff, Jitter: req.Jitter})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, x)
}

func (s *Server) deleteRetryPolicy(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteRetryPolicy(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type transitionRetryPolicyRequest struct {
	Status string `json:"status"`
}

func (s *Server) transitionRetryPolicy(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req transitionRetryPolicyRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	x, err := s.svc.TransitionRetryPolicy(id, req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, x)
}
