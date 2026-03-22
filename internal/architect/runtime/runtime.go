package runtime

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"project-generator/internal/architect/model"
	"project-generator/internal/config"
)

const ContextRelPath = ".architect/context.json"

func ResolveRoot(parentPath, projectName string) (string, error) {
	base, err := config.ExpandPath(parentPath)
	if err != nil {
		return "", err
	}
	if base == "" {
		base, err = config.ExpandPath(".")
		if err != nil {
			return "", err
		}
	}
	if projectName == "" {
		return base, nil
	}
	return filepath.Join(base, projectName), nil
}

func BackupIfNeeded(projectRoot string, overwrite bool) (string, error) {
	info, err := os.Stat(projectRoot)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("target path exists and is not a directory: %s", projectRoot)
	}
	if !overwrite {
		return "", fmt.Errorf("target path already exists: %s (use --overwrite)", projectRoot)
	}
	backupPath := fmt.Sprintf("%s_backup_%s", projectRoot, time.Now().Format("20060102150405"))
	if err := os.Rename(projectRoot, backupPath); err != nil {
		return "", err
	}
	return backupPath, nil
}

func SaveContext(projectRoot string, ctx model.ProjectContext) error {
	ctx.EnsureDefaults()
	contextPath := filepath.Join(projectRoot, ContextRelPath)
	if err := os.MkdirAll(filepath.Dir(contextPath), 0o755); err != nil {
		return err
	}
	payload, err := json.MarshalIndent(ctx, "", "  ")
	if err != nil {
		return err
	}
	payload = append(payload, '\n')
	return os.WriteFile(contextPath, payload, 0o644)
}

func LoadContext(projectRoot string) (model.ProjectContext, error) {
	contextPath := filepath.Join(projectRoot, ContextRelPath)
	raw, err := os.ReadFile(contextPath)
	if err != nil {
		return model.ProjectContext{}, err
	}
	var ctx model.ProjectContext
	if err := json.Unmarshal(raw, &ctx); err != nil {
		return model.ProjectContext{}, err
	}
	ctx.EnsureDefaults()
	return ctx, nil
}
