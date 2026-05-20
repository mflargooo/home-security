package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/mflargooo/internal/models"
	"github.com/mflargooo/internal/service"
	"github.com/mflargooo/internal/store"
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

// if rtsp url is given, read stream into mediamtx. otherwise assumes camera will push to expected path
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

	camera, result, err := h.svc.Register(r.Context(), req)
	if err != nil {
		log.Printf("register camera: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to register camera")
		return
	}

	httpStatus := http.StatusOK
	if result == store.UpsertCreated {
		httpStatus = http.StatusCreated
	}
	writeJSON(w, httpStatus, camera.ID)

	if httpStatus == http.StatusCreated {
		log.Printf("[REGISTER] created %s", camera.ID)
	} else {
		log.Printf("[REGISTER] upserted %s", camera.ID)
	}
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
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing camera id")
		return
	}

	camera, err := h.svc.Get(r.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "camera not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get camera")
		return
	}

	writeJSON(w, http.StatusOK, camera)
}

// if rtsp url is given, read stream into mediamtx. otherwise if rtsp url is cleared assumes camera will push to expected path
func (h *CameraHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing camera id")
		return
	}

	var req models.UpdateCameraRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	camera, err := h.svc.Update(r.Context(), id, req)
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "camera not found")
	}

	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update camera")
		return
	}

	writeJSON(w, http.StatusOK, camera)
}

func (h *CameraHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing camera id")
		return
	}

	err := h.svc.Delete(r.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "camera not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete camera")
		return
	}

	w.WriteHeader(http.StatusNoContent)
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
