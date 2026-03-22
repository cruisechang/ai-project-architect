package fsutil

import (
	"fmt"
	"os"
)

func EnsureDir(path string, dryRun bool) error {
	if dryRun {
		return nil
	}
	if err := os.MkdirAll(path, 0o755); err != nil {
		return fmt.Errorf("cannot create directory %s: %w", path, err)
	}
	return nil
}

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
