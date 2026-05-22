package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func MoveFile(src string, dst string) error { // deletes from src
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("mkdir clip path: %w", err)
	}

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer out.Close()

	log.Printf("[MOVING] from=%s to=%s", src, dst)
	if _, err := io.Copy(out, in); err != nil {
		os.Remove(dst)
		return fmt.Errorf("copy file: %w", err)
	}

	if err := os.Remove(src); err != nil {
		return fmt.Errorf("remove file: %v", err.Error())
	}

	return nil
}
