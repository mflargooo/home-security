package workers

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/mflargooo/internal/models"
	"github.com/mflargooo/internal/utils"
)

type ClipJob struct {
	ClipID    string `json:"id"`
	CameraID  string `json:"camera_id"`
	SessionID string
	Start     time.Time `json:"start_time"`
	End       time.Time `json:"end_time"`
	TmpPath   string
	FinalPath string `json:"file_path,omitempty"`
	Segments  []string
}

type ClipManager struct {
	resolver models.PathResolver
}

func NewClipManager(resolver models.PathResolver) *ClipManager {
	return &ClipManager{
		resolver: resolver,
	}
}

func (m *ClipManager) Dispatch(ctx context.Context, onDone func(clipID string, status models.ClipStatus, filePath string), jobs ...ClipJob) {
	for _, job := range jobs {
		go func(j ClipJob) {
			if err := utils.Retry(ctx, func() error { return runClipWorker(ctx, j) }, 3); err != nil {
				onDone(j.ClipID, "failed", "")
			} else {
				onDone(j.ClipID, "ready", j.FinalPath)
			}
		}(job)
	}
}

func writeConcatFile(path string, segments []string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	for _, seg := range segments {
		if _, err := fmt.Fprintf(w, "file '%s'\n", seg); err != nil {
			return err
		}
	}

	return w.Flush()
}

func runClipWorker(ctx context.Context, job ClipJob) error {
	if len(job.Segments) <= 0 {
		return fmt.Errorf("no segments")
	}

	if err := os.Remove(job.TmpPath); err != nil && !os.IsNotExist(err) {
		log.Printf("failed to remove tmp: %v", err)
	}

	sort.Strings(job.Segments)
	firstBase := filepath.Base(job.Segments[0])
	firstStem := strings.TrimSuffix(strings.TrimPrefix(firstBase, "seg-"), filepath.Ext(firstBase))

	firstStart, err := time.Parse("2006-01-02_15-04-05", firstStem)
	if err != nil {
		return err
	}

	ss := job.Start.Sub(firstStart).Seconds()
	duration := job.End.Sub(job.Start).Seconds()

	if duration <= 0 {
		return fmt.Errorf("invalid clip duration")
	}

	concatFile, _ := os.CreateTemp(os.TempDir(), "clip_*_concat.txt")
	concatPath := concatFile.Name()
	concatFile.Close()
	defer os.Remove(concatPath)

	if err := writeConcatFile(concatPath, job.Segments); err != nil {
		return err
	}

	cmd := exec.CommandContext(
		ctx,
		"ffmpeg",
		"-y",
		"-f", "concat",
		"-safe", "0",
		"-i", concatPath,
		"-ss", strconv.FormatFloat(ss, 'f', 3, 64),
		"-t", strconv.FormatFloat(duration, 'f', 3, 64),
		"-avoid_negative_ts", "make_zero",
		"-c", "copy",
		job.TmpPath,
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		log.Printf("ffmpeg failed: %v\n%s", err, stderr.String())
		return err
	}

	if err := utils.MoveFile(job.TmpPath, job.FinalPath); err != nil {
		return err
	}

	os.RemoveAll(filepath.Dir(job.Segments[0]))
	log.Printf("[CLIP] %s saved to %s", job.ClipID, job.FinalPath)

	return nil
}
