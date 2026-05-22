package service

import (
	"github.com/mflargooo/internal/models"
	"github.com/mflargooo/internal/store"
)

type ClipService struct {
	clipStore *store.ClipStore
	resolver  models.PathResolver
}

func NewClipService(clipStore *store.ClipStore, resolver models.PathResolver) *ClipService {
	return &ClipService{
		clipStore: clipStore,
		resolver:  resolver,
	}
}
