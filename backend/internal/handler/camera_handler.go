package handler

import (
	"net/http"

	"github.com/mflargooo/internal/service"
)

type CameraHandler struct {
	svc *service.CameraService
}

func NewCameraHandler(svc *service.CameraService) *CameraHandler {
	return &CameraHandler{svc: svc}
}

func (h *CameraHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"online":"true"}`))
	})

	mux.HandleFunc("POST /cameras", h.Register)
	mux.HandleFunc("GET /cameras", h.List)
	mux.HandleFunc("GET /cameras/{id}", h.GetByID)
	mux.HandleFunc("PATCH /cameras/{id}", h.Update)
	mux.HandleFunc("DELETE /cameras/{id}", h.Delete)
	mux.HandleFunc("POST /cameras/{id}/online", h.ReportOnline)
	mux.HandleFunc("POST /cameras/{id}/offline", h.ReportOffline)
}

func (h *CameraHandler) Register(w http.ResponseWriter, r *http.Request) {

}

func (h *CameraHandler) List(w http.ResponseWriter, r *http.Request) {

}

func (h *CameraHandler) GetByID(w http.ResponseWriter, r *http.Request) {

}

func (h *CameraHandler) Update(w http.ResponseWriter, r *http.Request) {

}

func (h *CameraHandler) Delete(w http.ResponseWriter, r *http.Request) {

}

func (h *CameraHandler) ReportOnline(w http.ResponseWriter, r *http.Request) {

}

func (h *CameraHandler) ReportOffline(w http.ResponseWriter, r *http.Request) {

}
