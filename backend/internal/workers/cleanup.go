package workers

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mflargooo/internal/models"
)

func StartCleanupWorker(resolver models.PathResolver) {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		log.Printf("[CLEANUP] started cleaner")
		for range ticker.C {
			cleanupSegments(resolver.BufferLivePath(""), 24*time.Hour)
		}
	}()
}

func cleanupSegments(root string, maxAge time.Duration) {
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
