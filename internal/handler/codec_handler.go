package handler

import (
	"net/http"

	"lightweightrpc/internal/model"
	"lightweightrpc/pkg/httpx"
)

func (s *Server) registerCodecRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/codecs", s.createCodec)
	mux.HandleFunc("GET /api/codecs", s.listCodec)
	mux.HandleFunc("GET /api/codecs/{id}", s.getCodec)
	mux.HandleFunc("PUT /api/codecs/{id}", s.updateCodec)
	mux.HandleFunc("DELETE /api/codecs/{id}", s.deleteCodec)
	mux.HandleFunc("PATCH /api/codecs/{id}/status", s.transitionCodec)
}

type createCodecRequest struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Version string `json:"version"`
	Schema  string `json:"schema"`
}

func (s *Server) createCodec(w http.ResponseWriter, r *http.Request) {
	var req createCodecRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	x, err := s.svc.CreateCodec(model.Codec{Name: req.Name, Type: req.Type, Version: req.Version, Schema: req.Schema})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, x)
}

func (s *Server) listCodec(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.CodecFilter{
		Type:   r.URL.Query().Get("type"),
		Status: r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListCodecs(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getCodec(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	x, err := s.svc.GetCodec(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, x)
}

type updateCodecRequest struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Version string `json:"version"`
	Schema  string `json:"schema"`
}

func (s *Server) updateCodec(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateCodecRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	x, err := s.svc.UpdateCodec(id, model.Codec{Name: req.Name, Type: req.Type, Version: req.Version, Schema: req.Schema})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, x)
}

func (s *Server) deleteCodec(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteCodec(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type transitionCodecRequest struct {
	Status string `json:"status"`
}

func (s *Server) transitionCodec(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req transitionCodecRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	x, err := s.svc.TransitionCodec(id, req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, x)
}
