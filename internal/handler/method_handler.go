package handler

import (
	"net/http"

	"lightweightrpc/internal/model"
	"lightweightrpc/pkg/httpx"
)

func (s *Server) registerMethodRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/methods", s.createMethod)
	mux.HandleFunc("GET /api/methods", s.listMethod)
	mux.HandleFunc("GET /api/methods/{id}", s.getMethod)
	mux.HandleFunc("PUT /api/methods/{id}", s.updateMethod)
	mux.HandleFunc("DELETE /api/methods/{id}", s.deleteMethod)
	mux.HandleFunc("PATCH /api/methods/{id}/status", s.transitionMethod)
}

type createMethodRequest struct {
	ServiceID  string `json:"service_id"`
	Name       string `json:"name"`
	InputType  string `json:"input_type"`
	OutputType string `json:"output_type"`
	TimeoutMs  int    `json:"timeout_ms"`
}

func (s *Server) createMethod(w http.ResponseWriter, r *http.Request) {
	var req createMethodRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	x, err := s.svc.CreateMethod(model.Method{ServiceID: req.ServiceID, Name: req.Name, InputType: req.InputType, OutputType: req.OutputType, TimeoutMs: req.TimeoutMs})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, x)
}

func (s *Server) listMethod(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.MethodFilter{
		ServiceId: r.URL.Query().Get("service_id"),
		Status:    r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListMethods(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getMethod(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	x, err := s.svc.GetMethod(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, x)
}

type updateMethodRequest struct {
	ServiceID  string `json:"service_id"`
	Name       string `json:"name"`
	InputType  string `json:"input_type"`
	OutputType string `json:"output_type"`
	TimeoutMs  int    `json:"timeout_ms"`
}

func (s *Server) updateMethod(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateMethodRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	x, err := s.svc.UpdateMethod(id, model.Method{ServiceID: req.ServiceID, Name: req.Name, InputType: req.InputType, OutputType: req.OutputType, TimeoutMs: req.TimeoutMs})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, x)
}

func (s *Server) deleteMethod(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteMethod(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type transitionMethodRequest struct {
	Status string `json:"status"`
}

func (s *Server) transitionMethod(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req transitionMethodRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	x, err := s.svc.TransitionMethod(id, req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, x)
}
