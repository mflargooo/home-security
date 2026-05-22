package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/mflargooo/internal/models"
	"github.com/mflargooo/internal/service"
)

type SnapshotHandler struct {
	svc *service.SnapshotService
}

func NewSnapshotHandler(svc *service.SnapshotService) *SnapshotHandler {
	return &SnapshotHandler{svc: svc}
}

func (h *SnapshotHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /snapshots/{cameraID}", h.StartSession) // response contains session id
	mux.HandleFunc("DELETE /snapshots/{sessionID}", h.EndSession)
	mux.HandleFunc("GET /snapshots/{sessionID}/index.m3u8", h.GetPlaylist)
	mux.HandleFunc("GET /snapshots/{sessionID}/{filename}", h.GetSnapshot)
}

func (h *SnapshotHandler) StartSession(w http.ResponseWriter, r *http.Request) {
	cameraID := r.PathValue("cameraID")
	cameraID = strings.TrimSpace(cameraID)

	if cameraID == "" {
		writeError(w, http.StatusBadRequest, "missing camera id")
		return
	}

	if _, err := uuid.Parse(cameraID); err != nil {
		writeError(w, http.StatusBadRequest, "invalid camera id")
		return
	}

	sessionID := uuid.New().String()

	if err := h.svc.CreateSnapshot(sessionID, cameraID); err != nil {
		log.Printf("create snapshot: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to create snapshot")
		return
	}

	if err := h.svc.StoreSession(r.Context(), sessionID, cameraID); err != nil {
		log.Printf("store session: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to store session")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.CreateSnapshotResponse{
		SessionID:   sessionID,
		PlaylistURL: fmt.Sprintf("/snapshots/%s/index.m3u8", sessionID),
	})
}

func (h *SnapshotHandler) EndSession(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("sessionID")

	if err := h.svc.RemoveSession(r.Context(), sessionID); err != nil {
		log.Printf("end session: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to end session")
		return
	}
}

func (h *SnapshotHandler) GetPlaylist(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("sessionID")

	playlist, err := h.svc.GeneratePlaylist(sessionID)
	if err != nil {
		log.Printf("generate playlist: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to start snapshot session")
		return
	}

	w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	w.Write([]byte(playlist))
}

func (h *SnapshotHandler) GetSnapshot(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("sessionID")
	filename := r.PathValue("filename")

	log.Printf("[SNAPSHOT] serving %s", h.svc.BuildSegmentPath(sessionID, filename))
	http.ServeFile(w, r, h.svc.BuildSegmentPath(sessionID, filename))
}
