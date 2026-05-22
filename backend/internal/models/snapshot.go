package models

type CreateSnapshotRequest struct {
	CameraID string `json:"camera_id"`
}

type CreateSnapshotResponse struct {
	SessionID   string `json:"session_id"`
	PlaylistURL string `json:"playlist_url"`
}

type GetPlaylistRequest struct {
	SessionID string `json:"session_id"`
	CameraID  string `json:"camera_id"`
}

type GetSegmentRequest struct {
	SessionID string `json:"session_id"`
	Filename  string `json:"filename"`
}
