package handler

import (
	"net/http"

	"github.com/mflargooo/internal/service"
)

type ClipHandler struct {
	svc *service.ClipService
}

func (h *ClipHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /clip/{id}", h.GetByID)
	mux.HandleFunc("GET /clip", h.List)
	mux.HandleFunc("POST /clip/{id}", h.Create)
	mux.HandleFunc("DELETE /clip/{id}", h.Delete)
}

func NewClipHandler(svc *service.ClipService) *ClipHandler {
	return &ClipHandler{svc: svc}
}

func (h *ClipHandler) GetByID(w http.ResponseWriter, r *http.Request) {

}

func (h *ClipHandler) List(w http.ResponseWriter, r *http.Request) {

}

func (h *ClipHandler) Create(w http.ResponseWriter, r *http.Request) {

}

func (h *ClipHandler) Delete(w http.ResponseWriter, r *http.Request) {

}
