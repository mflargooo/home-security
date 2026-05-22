package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mflargooo/internal/models"
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
	var req models.CreateClipRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("create clip: %v", err)
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validateClipCreateRequest(req); err != nil {
		log.Printf("create clip: %v", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	clipID, err := h.svc.Save(r.Context(), req)
	if err != nil {
		log.Printf("create clip: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to save clip")
		return
	}

	httpStatus := http.StatusOK
	if true { // in preparation for idemptoency, maybe
		httpStatus = http.StatusCreated
	}
	writeJSON(w, httpStatus, clipID)

	if httpStatus == http.StatusCreated {
		log.Printf("[CLIP] created %s", clipID)
	}
}

func (h *ClipHandler) Delete(w http.ResponseWriter, r *http.Request) {

}

func validateClipCreateRequest(req models.CreateClipRequest) error {
	if req.End.After(req.Start) {
		return nil
	}
	return fmt.Errorf("end time is before start time")
}
