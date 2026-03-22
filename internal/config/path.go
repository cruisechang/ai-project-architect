package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ExpandPath(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", nil
	}

	if strings.HasPrefix(trimmed, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("cannot resolve home directory: %w", err)
		}
		switch {
		case trimmed == "~":
			trimmed = home
		case strings.HasPrefix(trimmed, "~/") || strings.HasPrefix(trimmed, "~\\"):
			trimmed = filepath.Join(home, trimmed[2:])
		default:
			return "", fmt.Errorf("unsupported home path format: %q", raw)
		}
	}

	expanded := os.ExpandEnv(trimmed)
	abs, err := filepath.Abs(expanded)
	if err != nil {
		return "", fmt.Errorf("invalid path %q: %w", raw, err)
	}
	return abs, nil
}
