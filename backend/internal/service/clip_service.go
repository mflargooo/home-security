package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mflargooo/internal/models"
	"github.com/mflargooo/internal/store"
	"github.com/mflargooo/internal/workers"
)

type ClipService struct {
	clipStore     *store.ClipStore
	snapshotStore *store.SnapshotStore
	manager       *workers.ClipManager
	resolver      models.PathResolver
}

func NewClipService(clipStore *store.ClipStore, snapshotStore *store.SnapshotStore, manager *workers.ClipManager, resolver models.PathResolver) *ClipService {
	return &ClipService{
		clipStore:     clipStore,
		snapshotStore: snapshotStore,
		manager:       manager,
		resolver:      resolver,
	}
}

func (s *ClipService) Save(ctx context.Context, req models.CreateClipRequest) (string, error) {
	sessionID := req.SessionID
	cameraID := req.CameraID
	start := req.Start
	end := req.End

	storedCameraID, err := s.snapshotStore.Get(ctx, sessionID)
	if err != nil {
		return "", fmt.Errorf("get session: %w", err)
	}

	if storedCameraID == "" {
		return "", ErrSessionNotFound
	}

	segments, err := s.getSegmentsInRange(sessionID, start, end)
	if err != nil {
		return "", fmt.Errorf("get segments in range: %w", err)
	}

	if len(segments) == 0 {
		return "", fmt.Errorf("no segments found for time range")
	}

	clip := models.ClipResponse{
		CameraID: cameraID,
		Start:    start,
		End:      end,
	}
	// log clip 'processing'
	clipID, err := s.clipStore.Insert(ctx, clip)
	if err != nil {
		return "", fmt.Errorf("insert clip: %w", err)
	}

	// hardlink and give to clip job so we can retry on failure
	linkedSegments, err := s.createHardlinks(sessionID, clipID, segments)
	if err != nil {
		return "", fmt.Errorf("create hardlinks: %w", err)
	}

	filename := start.UTC().Format("2006-01-02_15-04-05") + ".ts"
	job := workers.ClipJob{
		ClipID:    clipID,
		Start:     start,
		End:       end,
		TmpPath:   filepath.Join(s.resolver.TmpClipsPath(clipID), ".tmp-clip.ts"),
		FinalPath: filepath.Join(s.resolver.SavedClipsPath(cameraID, clipID), filename),
		Segments:  linkedSegments,
	}

	s.manager.Dispatch(ctx,
		func(clipID string, status models.ClipStatus, filePath string) {
			s.clipStore.UpdateStatus(ctx, clipID, status, filePath)
		}, job)

	return clipID, nil
}

func (s *ClipService) createHardlinks(sessionID string, clipID string, segments []string) ([]string, error) {
	srcDir := s.resolver.BufferSnapshotPath(sessionID)
	dstDir := s.resolver.TmpClipsPath(clipID)

	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return nil, fmt.Errorf("create clip seg dir: %w", err)
	}

	var linkedSegments []string
	for _, seg := range segments {
		src := filepath.Join(srcDir, seg)
		dst := filepath.Join(dstDir, "seg-"+seg)
		if err := os.Link(src, dst); err != nil {
			if os.IsNotExist(err) {
				continue
			}

			return nil, fmt.Errorf("hardlink segment: %w", err)
		}
		linkedSegments = append(linkedSegments, dst)
	}

	return linkedSegments, nil
}

func (s *ClipService) getSegmentsInRange(sessionID string, start time.Time, end time.Time) ([]string, error) {
	snapshotDir := s.resolver.BufferSnapshotPath(sessionID)

	entries, err := os.ReadDir(snapshotDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrSessionNotFound
		}

		return nil, fmt.Errorf("read snapshot dir: %w", err)
	}

	var segments []string
	for i, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".ts") {
			continue
		}

		segStart := parseTimestampFromFilename(entry.Name())
		var segEnd time.Time
		if i+1 < len(entries) {
			segEnd = parseTimestampFromFilename(entries[i+1].Name())
		} else {
			segEnd = segStart.Add(60 * time.Second)
		}

		if segEnd.After(start) && segStart.Before(end) { // overlapping [start, end]
			segments = append(segments, entry.Name())
		}
	}

	return segments, nil
}
