package service

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/mflargooo/internal/models"
	"github.com/mflargooo/internal/store"
	"github.com/mflargooo/internal/utils"
)

var (
	ErrSessionNotFound = errors.New("snapshot session not found or expired")
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

func (s *SnapshotService) StoreSession(ctx context.Context, sessionID string, cameraID string) error {
	return s.snapshotStore.Set(ctx, sessionID, cameraID)
}

func (s *SnapshotService) RemoveSession(ctx context.Context, sessionID string) error {
	return s.snapshotStore.Delete(ctx, sessionID)
}

func (s *SnapshotService) GeneratePlaylist(sessionID string) (string, error) {
	csvPath := filepath.Join(s.resolver.BufferSnapshotPath(sessionID), "/hardlinks/segments.csv")
	f, err := os.Open(csvPath)
	if err != nil {
		return "", fmt.Errorf("failed to open segment list: %v", err)
	}
	defer f.Close()

	type segInfo struct {
		name     string
		duration float64
	}
	var infos []segInfo
	var maxDuration float64

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ",")
		if len(parts) != 3 {
			continue
		}
		name := filepath.Base(parts[0])
		start, _ := strconv.ParseFloat(parts[1], 64)
		end, _ := strconv.ParseFloat(parts[2], 64)
		duration := end - start
		if duration > maxDuration {
			maxDuration = duration
		}
		infos = append(infos, segInfo{name: name, duration: duration})
	}

	var b strings.Builder
	b.WriteString("#EXTM3U\n")
	b.WriteString("#EXT-X-VERSION:3\n")
	b.WriteString(fmt.Sprintf("#EXT-X-TARGETDURATION:%d\n", int(math.Ceil(maxDuration))))
	b.WriteString("#EXT-X-PLAYLIST-TYPE:VOD\n")
	b.WriteString("#EXT-X-MEDIA-SEQUENCE:0\n\n")

	for _, info := range infos {
		t := parseTimestampFromFilename(info.name)
		b.WriteString(fmt.Sprintf("#EXT-X-PROGRAM-DATE-TIME:%s\n",
			t.UTC().Format(time.RFC3339)))
		b.WriteString(fmt.Sprintf("#EXTINF:%.3f,\n", info.duration))
		b.WriteString(info.name)
		b.WriteString("\n\n")
	}

	b.WriteString("#EXT-X-ENDLIST\n")

	outpath := filepath.Join(s.resolver.BufferSnapshotPath(sessionID), "/playlist.m3u8")
	return outpath, os.WriteFile(outpath, []byte(b.String()), 0644)
}

func (s *SnapshotService) BuildSegmentListPath(sessionID string) string {
	return filepath.Join(s.resolver.BufferSnapshotPath(sessionID), "/playlist.m3u8")
}

func (s *SnapshotService) BuildSegmentHardlinkPath(sessionID string, filename string) string {
	return fmt.Sprintf("%s/hardlinks/%s", s.resolver.BufferSnapshotPath(sessionID), filename)
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
	path := filepath.Join(s.resolver.BufferSnapshotPath(sessionID), "/hardlinks")

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

func (s *SnapshotService) CreateHardlinks(sessionID string, cameraID string) error {
	srcDir := s.resolver.BufferLivePath(cameraID)
	dstDir := filepath.Join(s.resolver.BufferSnapshotPath(sessionID), "/hardlinks")

	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return err
	}

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

	srcCSV := filepath.Join(srcDir, "segments.csv")
	dstCSV := filepath.Join(dstDir, "segments.csv")

	if err := utils.CopyFile(srcCSV, dstCSV); err != nil {
		return fmt.Errorf("failed to snapshot segment list: %v", err)
	}

	log.Printf("[HARDLINKS] created %s", dstDir)
	return nil
}

func (s *SnapshotService) RemoveHardlinks(sessionID string) error {
	if err := os.RemoveAll(filepath.Join(s.resolver.BufferSnapshotPath(sessionID), "/hardlinks")); err != nil {
		return fmt.Errorf("remove hardlinks: %w", err)
	}
	return nil
}
