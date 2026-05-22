package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mflargooo/internal/models"
)

type ClipStore struct {
	db *pgxpool.Pool
}

func NewClipStore(db *pgxpool.Pool) *ClipStore {
	return &ClipStore{db: db}
}

func (s *ClipStore) Insert(ctx context.Context, clip models.ClipResponse) (string, error) {
	err := s.db.QueryRow(ctx, `
		INSERT INTO clips (camera_id, start_time, end_time, status)
		VALUES ($1, $2, $3, 'processing')
		RETURNING id
	`, clip.CameraID, clip.Start, clip.End).Scan(&clip.ClipID)
	if err != nil {
		return "", fmt.Errorf("insert clip: %w", err)
	}

	return clip.ClipID, nil
}

func (s *ClipStore) UpdateStatus(ctx context.Context, clipID string, status models.ClipStatus, filePath string) error {
	_, err := s.db.Exec(ctx, `
		UPDATE clips SET status = $1, file_path = $2 WHERE id = $3
	`, status, filePath, clipID)
	return err
}
