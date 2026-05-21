package handler

import (
	"net/http"

	"github.com/mflargooo/internal/service"
)

type VodHandler struct {
	svc *service.VodService
}

func NewVodHandler(svc *service.VodService) *VodHandler {
	return &VodHandler{}
}

func (h *VodHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /vod", h.List)
	mux.HandleFunc("GET /vod/{id}", h.GetByID)
	mux.HandleFunc("POST /save", h.Save)
}

func (h *VodHandler) List(w http.ResponseWriter, r *http.Request) {

}

func (h *VodHandler) GetByID(w http.ResponseWriter, r *http.Request) {

}

func (h *VodHandler) Save(w http.ResponseWriter, r *http.Request) {

}
