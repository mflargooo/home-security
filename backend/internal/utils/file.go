package utils

import (
	"fmt"
	"io"
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

	if _, err := io.Copy(out, in); err != nil {
		os.Remove(dst)
		return fmt.Errorf("copy file: %w", err)
	}

	return fmt.Errorf("remove file: %v", os.Remove(src).Error())
}
