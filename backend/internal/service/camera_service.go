package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/mflargooo/internal/models"
	"github.com/mflargooo/internal/store"
	"github.com/mflargooo/internal/workers"
)

type CameraService struct {
	cameraStore *store.CameraStore
	stateStore  *store.CameraStateStore
	manager     *workers.PullStreamManager
}

func NewCameraService(cameraStore *store.CameraStore, stateStore *store.CameraStateStore, manager *workers.PullStreamManager) *CameraService {
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
	camera, err := s.cameraStore.GetByID(ctx, id)
	if errors.Is(err, store.ErrNotFound) {
		return nil, store.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get camera: %w", err)
	}

	return &models.CameraResponse{Camera: *camera}, nil
}

func (s *CameraService) Update(ctx context.Context, id string, req models.UpdateCameraRequest) (*models.CameraResponse, error) {
	camera, err := s.cameraStore.Update(ctx, id, req)
	if errors.Is(err, store.ErrNotFound) {
		return nil, store.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("update camera: %w", err)
	}

	if camera.RTSPUrl == "" {
		s.manager.DeleteWorker(camera.ID)
	} else {
		s.onCameraUpdated(ctx, camera.ID, camera.RTSPUrl)
	}

	return &models.CameraResponse{Camera: *camera}, nil
}

func (s *CameraService) Delete(ctx context.Context, id string) error {
	if err := s.cameraStore.Delete(ctx, id); err != nil {
		return err
	}

	s.manager.DeleteWorker(id)

	return nil
}

func (s *CameraService) onCameraUpdated(ctx context.Context, id string, url string) {
	if url == "" {
		return
	}

	s.manager.UpsertWorker(ctx, id, url)
}
