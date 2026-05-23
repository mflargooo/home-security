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

type PullStreamWorker struct {
	id     string
	url    string
	cancel context.CancelFunc
}

type PullStreamManager struct {
	mu       sync.Mutex
	workers  map[string]*PullStreamWorker
	resolver models.PathResolver
}

func NewPullStreamManager(resolver models.PathResolver) *PullStreamManager {
	return &PullStreamManager{
		workers:  make(map[string]*PullStreamWorker),
		resolver: resolver,
	}
}

func (m *PullStreamManager) UpsertWorker(id string, url string) {
	m.mu.Lock()

	if old, ok := m.workers[id]; ok {
		if url == old.url {
			m.mu.Unlock()
			return
		}
		old.cancel()
	}

	ctx, cancel := context.WithCancel(context.Background())

	worker := &PullStreamWorker{
		id:     id,
		url:    url,
		cancel: cancel,
	}

	m.workers[id] = worker

	m.mu.Unlock()

	go m.runPullStreamWorker(ctx, id, url)
}

func (m *PullStreamManager) DeleteWorker(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if worker, ok := m.workers[id]; ok {
		worker.cancel()
		delete(m.workers, id)
	}
}

func (m *PullStreamManager) runPullStreamWorker(ctx context.Context, id string, url string) {
	streamLiveURL := m.resolver.StreamLiveURL(id)
	bufferLivePath := m.resolver.BufferLivePath(id)

	if err := os.MkdirAll(bufferLivePath, 0755); err != nil {
		log.Printf("[FFMPEG] buffer live path for worker=%s: %v\n", id, err)
		return
	}

	teeOutput := fmt.Sprintf(
		"[f=rtsp:rtsp_transport=tcp]%s|[f=segment:segment_time=4:strftime=1:reset_timestamps=1]%s/%%Y-%%m-%%d_%%H-%%M-%%S.ts",
		streamLiveURL,
		bufferLivePath,
	)

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
