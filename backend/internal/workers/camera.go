package workers

import (
	"bytes"
	"context"
	"log"
	"os/exec"
	"sync"
	"time"

	"github.com/mflargooo/internal/models"
)

type FFmpegWorker struct {
	id     string
	url    string
	cancel context.CancelFunc
}

type FFmpegManager struct {
	mu       sync.Mutex
	workers  map[string]*FFmpegWorker
	resolver models.StreamResolver
}

func NewFFmpegManager(resolver models.StreamResolver) *FFmpegManager {
	return &FFmpegManager{
		workers:  make(map[string]*FFmpegWorker),
		resolver: resolver,
	}
}

func (m *FFmpegManager) UpsertWorker(id string, url string) {
	m.mu.Lock()

	if old, ok := m.workers[id]; ok {
		if url == old.url {
			m.mu.Unlock()
			return
		}
		old.cancel()
	}

	ctx, cancel := context.WithCancel(context.Background())

	worker := &FFmpegWorker{
		id:     id,
		url:    url,
		cancel: cancel,
	}

	m.workers[id] = worker

	m.mu.Unlock()

	go m.runFFmpegWorker(ctx, id, url)
}

func (m *FFmpegManager) DeleteWorker(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if worker, ok := m.workers[id]; ok {
		worker.cancel()
		delete(m.workers, id)
	}
}

func (m *FFmpegManager) runFFmpegWorker(ctx context.Context, id string, url string) {
	outputURL := m.resolver.StreamURL("/cameras/" + id + "/stream")
	log.Printf("[FFMPEG] worker=%s streaming video to public=%q", id, outputURL)

	for {
		if ctx.Err() != nil {
			return
		}

		cmd := exec.CommandContext(
			ctx,
			"ffmpeg",
			"-rtsp_transport", "tcp",
			"-i", url,
			"-c", "copy",
			"-f", "rtsp",
			"-rtsp_transport", "tcp",
			outputURL,
		)

		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		err := cmd.Run()

		if ctx.Err() != nil {
			return
		}

		log.Printf("[FFMPEG] worker=%s crashed: %v", id, err)
		log.Printf("[FFMPEG] stderr: %s\n", stderr.String())

		select {
		case <-ctx.Done():
			return
		case <-time.After(5 * time.Second):
		}
	}
}
