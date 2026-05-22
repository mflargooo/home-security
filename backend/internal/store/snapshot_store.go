package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type SnapshotStore struct {
	client *redis.Client
}

func NewSnapshotStore(client *redis.Client) *SnapshotStore {
	return &SnapshotStore{client: client}
}

func sessionKey(sessionID string) string {
	return "snapshot:session:" + sessionID
}

func (s *SnapshotStore) Set(ctx context.Context, sessionID string, cameraID string) error {
	return s.client.Set(ctx, sessionKey(sessionID), cameraID, 2*time.Hour).Err()
}

func (s *SnapshotStore) Get(ctx context.Context, sessionID string) (string, error) {
	cameraID, err := s.client.Get(ctx, sessionKey(sessionID)).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("redis get session: %w", err)
	}

	return cameraID, nil
}

func (s *SnapshotStore) Exists(ctx context.Context, sessionID string) (int64, error) {
	return s.client.Exists(ctx, sessionKey(sessionID)).Result()
}

func (s *SnapshotStore) Delete(ctx context.Context, sessionID string) error {
	return s.client.Del(ctx, sessionKey(sessionID)).Err()
}
