package store

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("camera not found")

type CameraStore struct {
	db *pgxpool.Pool
}

func NewCameraStore(db *pgxpool.Pool) *CameraStore {
	return &CameraStore{db: db}
}
