package handler

import (
	"net/http"

	"lightweightrpc/internal/model"
	"lightweightrpc/pkg/httpx"
)

func (s *Server) registerCallLogRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/call-logs", s.createCallLog)
	mux.HandleFunc("GET /api/call-logs", s.listCallLog)
	mux.HandleFunc("GET /api/call-logs/{id}", s.getCallLog)
	mux.HandleFunc("DELETE /api/call-logs/{id}", s.deleteCallLog)
}

type createCallLogRequest struct {
	ServiceID string `json:"service_id"`
	MethodID  string `json:"method_id"`
	Caller    string `json:"caller"`
	LatencyMs int64  `json:"latency_ms"`
	Code      int    `json:"code"`
	ErrMsg    string `json:"err_msg"`
	Timestamp int64  `json:"timestamp"`
}

func (s *Server) createCallLog(w http.ResponseWriter, r *http.Request) {
	var req createCallLogRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	x, err := s.svc.CreateCallLog(model.CallLog{ServiceID: req.ServiceID, MethodID: req.MethodID, Caller: req.Caller, LatencyMs: req.LatencyMs, Code: req.Code, ErrMsg: req.ErrMsg, Timestamp: req.Timestamp})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, x)
}

func (s *Server) listCallLog(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.CallLogFilter{
		ServiceId: r.URL.Query().Get("service_id"),
		Caller:    r.URL.Query().Get("caller"),
	}
	items, total, err := s.svc.ListCallLogs(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getCallLog(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	x, err := s.svc.GetCallLog(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, x)
}

func (s *Server) deleteCallLog(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteCallLog(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
