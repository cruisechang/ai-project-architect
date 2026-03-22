package fsutil

import (
	"fmt"
	"os"
	"path/filepath"
)

func WriteFile(path string, data []byte, perm os.FileMode, dryRun bool) error {
	if dryRun {
		return nil
	}
	if perm == 0 {
		perm = 0o644
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("cannot create parent directory for %s: %w", path, err)
	}
	if err := os.WriteFile(path, data, perm); err != nil {
		return fmt.Errorf("cannot write file %s: %w", path, err)
	}
	return nil
}
