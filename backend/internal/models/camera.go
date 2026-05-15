package models

import "time"

type CameraStatus string

const (
	StatusActive       CameraStatus = "active"
	StatusOffline      CameraStatus = "offline"
	StatusDisabled     CameraStatus = "disabled"
	StatusReconnecting CameraStatus = "reconnecting"
)

type Camera struct {
	ID         string            `json:"id" db:"id"`
	Name       string            `json:"name" db:"name"`
	RTSPUrl    string            `json:"rtsp_url" db:"rtsp_url"`
	MACAddress string            `json:"mac_address" db:"mac_address"`
	Status     CameraStatus      `json:"status" db:"status"`
	Metadata   map[string]string `json:"metadata,omitempty" db:"metadata"`
	CreatedAt  time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at" db:"updated_at"`
	LastSeenAt *time.Time        `json:"last_seen_at,omitempty" db:"last_seen_at"`
}

type CameraState struct {
	CameraID        string       `json:"camera_id"`
	Status          CameraStatus `json:"status"`
	ReconnectCount  int          `json:"reconnected_count"`
	LastReconnectAt *time.Time   `json:"last_reconnect_at,omitempty"`
	LastError       string       `json:"last_error,omitempty"`
	NextReconnectAt *time.Time   `json:"next_reconnected_at,omitempty"`
}

type CreateCameraRequest struct {
	Name       string            `json:"name" validate:"required,min=1,max=100"`
	RTSPUrl    string            `json:"rtsp_url" validate:"required,url"`
	MACAddress string            `json:"mac_address" validate:"omitempty,mac"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

type UpdateCameraRequest struct {
	Name       *string           `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	RTSPUrl    *string           `json:"rtsp_url,omitempty" validate:"omitempty,url"`
	Status     *CameraStatus     `json:"status,omitempty"`
	MACAddress *string           `json:"mac_address,omitempty" validate:"omitempty,mac"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

type CameraResponse struct {
	Camera
	State *CameraState `json:"state,omitempty"`
}
