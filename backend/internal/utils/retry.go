package utils

import (
	"context"
	"fmt"
	"log"
	"time"
)

func Retry(ctx context.Context, f func() error, maxAttempts int) error {
	var lastErr error
	for attempt := 0; attempt < maxAttempts; attempt++ {
		log.Printf("[RETRY] attempt %d/%d", attempt+1, maxAttempts)
		lastErr = f()
		if lastErr == nil {
			return nil
		}
		log.Printf("[RETRY] attempt %d failed: %v", attempt+1, lastErr)
		if attempt < maxAttempts-1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(time.Duration((attempt+1)*5) * time.Second):
			}
		}
	}
	return fmt.Errorf("failed after %d attempts: %w", maxAttempts, lastErr)
}
