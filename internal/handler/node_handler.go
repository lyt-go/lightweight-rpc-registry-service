package handler

import (
	"net/http"

	"lightweightrpc/internal/model"
	"lightweightrpc/pkg/httpx"
)

func (s *Server) registerNodeRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/nodes", s.createNode)
	mux.HandleFunc("GET /api/nodes", s.listNode)
	mux.HandleFunc("GET /api/nodes/{id}", s.getNode)
	mux.HandleFunc("PUT /api/nodes/{id}", s.updateNode)
	mux.HandleFunc("DELETE /api/nodes/{id}", s.deleteNode)
	mux.HandleFunc("PATCH /api/nodes/{id}/status", s.transitionNode)
}

type createNodeRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Weight  int    `json:"weight"`
}

func (s *Server) createNode(w http.ResponseWriter, r *http.Request) {
	var req createNodeRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	x, err := s.svc.CreateNode(model.Node{Name: req.Name, Address: req.Address, Region: req.Region, Weight: req.Weight})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, x)
}

func (s *Server) listNode(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.NodeFilter{
		Region: r.URL.Query().Get("region"),
		Status: r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListNodes(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getNode(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	x, err := s.svc.GetNode(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, x)
}

type updateNodeRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Weight  int    `json:"weight"`
}

func (s *Server) updateNode(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateNodeRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	x, err := s.svc.UpdateNode(id, model.Node{Name: req.Name, Address: req.Address, Region: req.Region, Weight: req.Weight})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, x)
}

func (s *Server) deleteNode(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteNode(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type transitionNodeRequest struct {
	Status string `json:"status"`
}

func (s *Server) transitionNode(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req transitionNodeRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	x, err := s.svc.TransitionNode(id, req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, x)
}
