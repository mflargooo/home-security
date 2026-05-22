package workers

import (
	"context"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mflargooo/internal/models"
	"github.com/mflargooo/internal/store"
	"github.com/mflargooo/internal/utils"
)

func StartBufferCleanupWorker(resolver models.PathResolver, maxAge time.Duration) {
	utils.Interval(func() {
		cleanupBufferSegments(resolver.BufferLivePath(""), maxAge)
	}, 5*time.Minute)
	log.Printf("[CLEANUP] started buffer cleaner")
}

func cleanupBufferSegments(root string, maxAge time.Duration) {
	root = filepath.Clean(root)
	cutoff := time.Now().Add(-maxAge)

	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".mp4") {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}

		if info.ModTime().Before(cutoff) {
			log.Printf("[CLEANUP] removing: %s", path)
			os.Remove(path)
		}

		return nil
	})
}

func StartSnapshotCleanupWorker(resolver models.PathResolver, snapshotStore *store.SnapshotStore) {
	utils.Interval(func() {
		cleanupSnapshotSegments(resolver.BufferSnapshotPath(""), snapshotStore)
	}, 30*time.Minute)
	log.Printf("[CLEANUP] started snapshot cleaner")
}

func cleanupSnapshotSegments(root string, snapshotStore *store.SnapshotStore) {
	entries, err := os.ReadDir(root)
	if err != nil {
		log.Printf("[CLEANUP] failed to read snapshot dir: %v", err)
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		sessionUUID, err := uuid.Parse(entry.Name())
		sessionID := sessionUUID.String()
		if err != nil {
			log.Printf("[CLEANUP] snapshot found invalid session id %s: %v", sessionID, err)
			continue
		}

		exists, err := snapshotStore.Exists(context.Background(), sessionID)
		if err != nil {
			log.Printf("[CLEANUP] failed to check redis for %s: %v", sessionID, err)
			continue
		}

		if exists == 0 {
			path := root + "/" + entry.Name()
			log.Printf("[CLEANUP] removing expired snapshot %s", path)
			os.RemoveAll(path)
		}
	}
}
