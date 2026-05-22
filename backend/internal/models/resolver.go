package models

type PathResolver interface {
	StreamLiveURL(cameraID string) string
	StreamSnapshotURL(sessionID string, segment string) string
	BufferLivePath(cameraID string) string
	BufferSnapshotPath(sessionID string) string
}
