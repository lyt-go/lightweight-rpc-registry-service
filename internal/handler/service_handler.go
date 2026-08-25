package handler

import (
	"net/http"

	"lightweightrpc/internal/model"
	"lightweightrpc/pkg/httpx"
)

func (s *Server) registerServiceRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/services", s.createService)
	mux.HandleFunc("GET /api/services", s.listService)
	mux.HandleFunc("GET /api/services/{id}", s.getService)
	mux.HandleFunc("PUT /api/services/{id}", s.updateService)
	mux.HandleFunc("DELETE /api/services/{id}", s.deleteService)
	mux.HandleFunc("PATCH /api/services/{id}/status", s.transitionService)
}

type createServiceRequest struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
}

func (s *Server) createService(w http.ResponseWriter, r *http.Request) {
	var req createServiceRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	x, err := s.svc.CreateService(model.Service{Name: req.Name, Version: req.Version, Host: req.Host, Port: req.Port, Protocol: req.Protocol})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, x)
}

func (s *Server) listService(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ServiceFilter{
		Protocol: r.URL.Query().Get("protocol"),
		Status:   r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListServices(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getService(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	x, err := s.svc.GetService(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, x)
}

type updateServiceRequest struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
}

func (s *Server) updateService(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateServiceRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	x, err := s.svc.UpdateService(id, model.Service{Name: req.Name, Version: req.Version, Host: req.Host, Port: req.Port, Protocol: req.Protocol})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, x)
}

func (s *Server) deleteService(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteService(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type transitionServiceRequest struct {
	Status string `json:"status"`
}

func (s *Server) transitionService(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req transitionServiceRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	x, err := s.svc.TransitionService(id, req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, x)
}
