package service

import (
	"context"

	"github.com/mflargooo/internal/models"
	"github.com/mflargooo/internal/store"
)

type CameraService struct {
	cameraStore *store.CameraStore
	stateStore  *store.CameraStateStore
}

func NewCameraService(cameraStore *store.CameraStore, stateStore *store.CameraStateStore) *CameraService {
	return &CameraService{
		cameraStore: cameraStore,
		stateStore:  stateStore,
	}
}

func (s *CameraService) Register(ctx context.Context, req models.CreateCameraRequest) {

}
