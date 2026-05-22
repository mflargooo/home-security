package workers

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
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
	resolver models.PathResolver
}

func NewFFmpegManager(resolver models.PathResolver) *FFmpegManager {
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
	streamLiveURL := m.resolver.StreamLiveURL(id)
	bufferLivePath := m.resolver.BufferLivePath(id)

	if err := os.MkdirAll(bufferLivePath, 0755); err != nil {
		log.Printf("[FFMPEG] buffer live path for worker=%s: %v\n", id, err)
		return
	}

	teeOutput := fmt.Sprintf("[f=rtsp:rtsp_transport=tcp]%s|[f=segment:segment_time=60:strftime=1:reset_timestamps=1:segment_format=mp4]%s/%%Y-%%m-%%d_%%H-%%M-%%S.mp4", streamLiveURL, bufferLivePath)

	log.Printf("[FFMPEG] worker=%s streaming=%q writing=%q", id, streamLiveURL, bufferLivePath)

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
			"-f", "tee",
			"-map", "0",
			teeOutput,
		)

		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		err := cmd.Run()

		if err != nil {
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
