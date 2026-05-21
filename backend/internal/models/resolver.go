package models

type PathResolver interface {
	StreamLiveURL(cameraID string) string
	BufferLivePath(cameraID string) string
	BufferSnapshotPath(cameraID string, sessionID string) string
}
