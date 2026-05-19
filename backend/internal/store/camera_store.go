package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mflargooo/internal/models"
)

var ErrNotFound = errors.New("camera not found")

type CameraStore struct {
	db *pgxpool.Pool
}

func NewCameraStore(db *pgxpool.Pool) *CameraStore {
	return &CameraStore{db: db}
}

type UpsertResult int

const (
	UpsertCreated  UpsertResult = iota // new cam
	UpsertReturned                     // registered url
)

// IDEMPOTENT: true if new insert, false if old returned
func (s *CameraStore) Create(ctx context.Context, req models.CreateCameraRequest) (*models.Camera, UpsertResult, error) {
	metaJSON, err := json.Marshal(req.Metadata)
	if err != nil {
		return nil, 0, fmt.Errorf("marshal metadata: %w", err)
	}

	now := time.Now().UTC()
	id := uuid.New().String()

	query := `
		INSERT INTO cameras (id, name, rtsp_url, status, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (rtsp_url) DO UPDATE SET
			name = EXCLUDED.name,
			metadata = EXCLUDED.metadata,
			updated_at = EXCLUDED.updated_at
		RETURNING id, name, rtsp_url, status, metadata, created_at, updated_at, (xmax = 0) AS inserted
	`

	var camera models.Camera
	var metaOut []byte
	var wasInserted bool

	err = s.db.QueryRow(ctx, query,
		id,
		req.Name,
		req.RTSPUrl,
		models.StatusActive,
		metaJSON,
		now,
		now,
	).Scan(
		&camera.ID,
		&camera.Name,
		&camera.RTSPUrl,
		&camera.Status,
		&metaOut,
		&now,
		&now,
		&wasInserted,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("upsert camera: %w", err)
	}

	if len(metaOut) > 0 {
		if err := json.Unmarshal(metaOut, &camera.Metadata); err != nil {
			return nil, 0, fmt.Errorf("unmarshal metadata: %w", err)
		}
	}

	var result UpsertResult
	if wasInserted {
		result = UpsertCreated
	} else {
		result = UpsertReturned
	}

	return &camera, result, nil
}

func (s *CameraStore) GetByID(ctx context.Context, id string) (*models.Camera, error) {
	query := `
		SELECT id, name, rtsp_url, status, metadata, created_at, updated_at, last_seen_at
		FROM cameras
		WHERE id = $1
	`

	var camera models.Camera
	var metaJSON []byte

	err := s.db.QueryRow(ctx, query, id).Scan(
		&camera.ID,
		&camera.Name,
		&camera.RTSPUrl,
		&camera.Status,
		&camera.Metadata,
		&camera.CreatedAt,
		&camera.UpdatedAt,
		&camera.LastSeenAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("query camera by id: %w", err)
	}

	if len(metaJSON) > 0 {
		if err := json.Unmarshal(metaJSON, &camera.Metadata); err != nil {
			return nil, fmt.Errorf("unmarshal metadata: %w", err)
		}
	}

	return &camera, nil
}

func (s *CameraStore) GetByRTSPURL(ctx context.Context, url string) (*models.Camera, error) {
	query := `
		SELECT id, name, rtsp_url, status, metadata, created_at, updated_at, last_seen_at
		FROM cameras
		WHERE rtsp_url = $1
	`

	var camera models.Camera
	var metaJSON []byte

	err := s.db.QueryRow(ctx, query, url).Scan(
		&camera.ID,
		&camera.Name,
		&camera.RTSPUrl,
		&camera.Status,
		&camera.Metadata,
		&camera.CreatedAt,
		&camera.UpdatedAt,
		&camera.LastSeenAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("query camera by rtsp_url: %w", err)
	}

	if len(metaJSON) > 0 {
		if err := json.Unmarshal(metaJSON, &camera.Metadata); err != nil {
			return nil, fmt.Errorf("unmarshal metadata: %w", err)
		}
	}

	return &camera, nil
}

func (s *CameraStore) List(ctx context.Context, status *models.CameraStatus) ([]*models.Camera, error) {
	query := `
		SELECT id, name, rtsp_url, status, metadata, created_at, updated_at, last_seen_at
		FROM cameras
	`

	args := []any{}

	if status != nil {
		query += "WHERE status = $1"
		args = append(args, *status)
	}

	query += "ORDER BY created_at DESC"

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list camers: %w", err)
	}
	defer rows.Close()

	var cameras []*models.Camera
	for rows.Next() {
		var camera models.Camera
		var metaJSON []byte

		err := rows.Scan(
			&camera.ID,
			&camera.Name,
			&camera.RTSPUrl,
			&camera.Status,
			&metaJSON,
			&camera.CreatedAt,
			&camera.UpdatedAt,
			&camera.LastSeenAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan camera row: %w", err)
		}

		if len(metaJSON) > 0 {
			if err := json.Unmarshal(metaJSON, &camera.Metadata); err != nil {
				return nil, fmt.Errorf("unmarshal metadata: %w", err)
			}
		}

		cameras = append(cameras, &camera)
	}

	return cameras, rows.Err()
}

func (s *CameraStore) Update(ctx context.Context, id string, req models.UpdateCameraRequest) (*models.Camera, error) {
	existing, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.RTSPUrl != nil {
		existing.RTSPUrl = *req.RTSPUrl
	}
	if req.Status != nil {
		existing.Status = *req.Status
	}
	if req.Metadata != nil {
		existing.Metadata = req.Metadata
	}

	existing.UpdatedAt = time.Now().UTC()

	metaJSON, err := json.Marshal(existing.Metadata)
	if err != nil {
		return nil, fmt.Errorf("marshal metadata: %w", err)
	}

	query := `
		UPDATE cameras
		SET name=$1, rtsp_url=$2, status=$3, metadata=$4, updated_at=$5
		WHERE id=$6
	`

	_, err = s.db.Exec(ctx, query,
		existing.Name,
		existing.RTSPUrl,
		existing.Status,
		metaJSON,
		existing.UpdatedAt,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("update camera: %w", err)
	}

	return existing, nil
}

func (s *CameraStore) UpdateLastSeen(ctx context.Context, id string) error {
	now := time.Now().UTC()
	_, err := s.db.Exec(
		ctx,
		"UPDATE cameras SET last_seen_at=$1 WHERE id=$2",
		now, id,
	)

	return err
}

func (s *CameraStore) Delete(ctx context.Context, id string) error {
	result, err := s.db.Exec(ctx, "DELETE FROM cameras WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("delete camera: %w", err)
	}

	rows := result.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}
