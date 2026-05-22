package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mflargooo/internal/models"
	"github.com/mflargooo/internal/store"
)

type SnapshotService struct {
	snapshotStore *store.SnapshotStore
	resolver      models.PathResolver
}

func NewSnapshotService(snapshotStore *store.SnapshotStore, resolver models.PathResolver) *SnapshotService {
	return &SnapshotService{
		snapshotStore: snapshotStore,
		resolver:      resolver,
	}
}

// Delete handled by dispatched worker checking Redis
func (s *SnapshotService) CreateSnapshot(sessionID string, cameraID string) error {
	return s.createHardlinks(sessionID, cameraID)
}

func (s *SnapshotService) StoreSession(ctx context.Context, sessionID string, cameraID string) error {
	return s.snapshotStore.Set(ctx, sessionID, cameraID)
}

func (s *SnapshotService) RemoveSession(ctx context.Context, sessionID string) error {
	return s.snapshotStore.Delete(ctx, sessionID)
}

func (s *SnapshotService) GeneratePlaylist(sessionID string) (string, error) {
	segments, err := s.getSegments(sessionID)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve segments: %v", err)
	}
	var b strings.Builder
	b.WriteString("#EXTM3U\n")
	b.WriteString("#EXT-X-VERSION:3\n")
	b.WriteString("#EXT-X-TARGETDURATION:60\n")
	b.WriteString("#EXT-X-PLAYLIST-TYPE:VOD\n\n")

	for _, seg := range segments {
		t := parseTimestampFromFilename(seg)
		b.WriteString(fmt.Sprintf("#EXT-X-PROGRAM-DATE-TIME:%s\n", t.UTC().Format(time.RFC3339)))
		b.WriteString("#EXTINF:60.0,\n")
		b.WriteString(s.resolver.StreamSnapshotURL(sessionID, seg))
		b.WriteString("\n\n")
	}

	b.WriteString("#EXT-X-ENDLIST\n")
	return b.String(), nil
}

func (s *SnapshotService) BuildSegmentPath(sessionID string, filename string) string {
	return fmt.Sprintf("%s/%s", s.resolver.BufferSnapshotPath(sessionID), filename)
}

func parseTimestampFromFilename(filename string) time.Time {
	name := strings.TrimSuffix(filename, ".ts")
	t, err := time.ParseInLocation("2006-01-02_15-04-05", name, time.UTC)
	if err != nil {
		return time.Time{}
	}
	return t
}

func (s *SnapshotService) getSegments(sessionID string) ([]string, error) {
	path := s.resolver.BufferSnapshotPath(sessionID)

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var segments []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".ts") {
			continue
		}

		segments = append(segments, entry.Name())
	}
	return segments, nil
}

func (s *SnapshotService) createHardlinks(sessionID string, cameraID string) error {
	srcDir := s.resolver.BufferLivePath(cameraID)
	dstDir := s.resolver.BufferSnapshotPath(sessionID)

	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return err
	}

	log.Printf("[HARDLINKS] created %s", dstDir)

	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}

	if len(entries) > 1 { // exclude the most recent segment being written to, so snapshot does not grow
		entries = entries[:len(entries)-1]
	}

	for _, entry := range entries {
		src := filepath.Join(srcDir, entry.Name())
		dst := filepath.Join(dstDir, entry.Name())

		if err := os.Link(src, dst); err != nil {
			if os.IsNotExist(err) {
				continue
			}

			return err
		}
	}

	return nil
}
