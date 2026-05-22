package models

import "time"

type ClipStatus string

const (
	StatusReady      ClipStatus = "ready"
	StatusProcessing ClipStatus = "processing"
	StatusFailed     ClipStatus = "failed"
)

type CreateClipRequest struct {
	SessionID string
	Start     time.Time `json:"start_time"`
	End       time.Time `json:"end_time"`
}

type ClipResponse struct {
	ClipID   string     `json:"id"`
	CameraID string     `json:"camera_id"`
	Start    time.Time  `json:"start_time"`
	End      time.Time  `json:"end_time"`
	FilePath string     `json:"file_path,omitempty"`
	Status   ClipStatus `json:"status"`
}
