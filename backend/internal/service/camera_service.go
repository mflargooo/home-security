package service

import (
	"context"
	"fmt"

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

func (s *CameraService) Register(ctx context.Context, req models.CreateCameraRequest) (*models.CameraResponse, error) {
	camera, err := s.cameraStore.Create(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("register camera: %w", err)
	}

	return &models.CameraResponse{Camera: *camera}, nil
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
