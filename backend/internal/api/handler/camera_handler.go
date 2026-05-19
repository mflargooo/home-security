package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mflargooo/internal/models"
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
	var req models.CreateCameraRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("register camera: %v", err)
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validateCreateRequest(req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	camera, err := h.svc.Register(r.Context(), req)
	if err != nil {
		log.Printf("register camera: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to register camera")
		return
	}

	writeJSON(w, http.StatusCreated, camera.ID)
	log.Printf("[CREATE] camera created with: %s", camera.ID)
}

func (h *CameraHandler) List(w http.ResponseWriter, r *http.Request) {
	var statusFilter *models.CameraStatus
	if s := r.URL.Query().Get("status"); s != "" {
		status := models.CameraStatus(s)
		statusFilter = &status
	}

	cameras, err := h.svc.List(r.Context(), statusFilter)
	if err != nil {
		log.Printf("list cameras: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to list cameras")
		return
	}

	writeJSON(w, http.StatusOK, cameras)
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

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func validateCreateRequest(req models.CreateCameraRequest) error {
	return nil
}
