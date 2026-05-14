package store

import (
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	cameraStatePrefix = "camera:state:"
	cameraStateTTL    = 24 * time.Hour // expire cameras that do not update within 24 hrs

	// backoff reconnection
	MaxReconnectCount  = 10
	BaseReconnectDelay = 5 * time.Second
	MaxReconnectDelay  = 5 * time.Minute
)

type CameraStateStore struct {
	client *redis.Client
}

func NewCameraStateStore(client *redis.Client) *CameraStateStore {
	return &CameraStateStore{client: client}
}
