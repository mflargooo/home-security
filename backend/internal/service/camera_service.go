package service

import (
	"context"
	"fmt"

	"github.com/mflargooo/internal/models"
	"github.com/mflargooo/internal/store"
	"github.com/mflargooo/internal/workers"
)

type CameraService struct {
	cameraStore *store.CameraStore
	stateStore  *store.CameraStateStore
	manager     *workers.FFmpegManager
}

func NewCameraService(cameraStore *store.CameraStore, stateStore *store.CameraStateStore, manager *workers.FFmpegManager) *CameraService {
	return &CameraService{
		cameraStore: cameraStore,
		stateStore:  stateStore,
		manager:     manager,
	}
}

func (s *CameraService) Register(ctx context.Context, req models.CreateCameraRequest) (*models.CameraResponse, store.UpsertResult, error) {
	camera, result, err := s.cameraStore.Create(ctx, req)
	if err != nil {
		return nil, 0, fmt.Errorf("register camera: %w", err)
	}

	s.onCameraUpdated(ctx, camera.ID, camera.RTSPUrl)

	return &models.CameraResponse{Camera: *camera}, result, nil
}

func (s *CameraService) List(ctx context.Context, status *models.CameraStatus) ([]*models.CameraResponse, error) {
	cameras, err := s.cameraStore.List(ctx, status)
	if err != nil {
		return nil, fmt.Errorf("list cameras: %w", err)
	}

	responses := make([]*models.CameraResponse, len(cameras))

	for i, c := range cameras {
		responses[i] = &models.CameraResponse{Camera: *c}
	}

	return responses, nil

}

func (s *CameraService) Get(ctx context.Context, id string) (*models.CameraResponse, error) {
	return nil, nil
}

func (s *CameraService) onCameraUpdated(ctx context.Context, id string, url string) {
	if url == "" {
		return
	}

	s.manager.UpsertWorker(id, url)
}
