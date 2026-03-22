package scaffold

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerate(t *testing.T) {
	root := filepath.Join(t.TempDir(), "project")

	result, err := Generate(root)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	if len(result.CreatedDirs) == 0 {
		t.Error("no directories created")
	}
	if len(result.CreatedFiles) == 0 {
		t.Error("no files created")
	}

	for _, d := range result.CreatedDirs {
		info, err := os.Stat(d)
		if err != nil {
			t.Errorf("dir %q not found: %v", d, err)
			continue
		}
		if !info.IsDir() {
			t.Errorf("%q is not a directory", d)
		}
	}

	for _, f := range result.CreatedFiles {
		info, err := os.Stat(f)
		if err != nil {
			t.Errorf("file %q not found: %v", f, err)
			continue
		}
		if info.IsDir() {
			t.Errorf("%q is a directory, expected file", f)
		}
		if info.Size() == 0 {
			t.Errorf("%q is empty", f)
		}
	}
}

func TestGenerateExpectedDirs(t *testing.T) {
	root := filepath.Join(t.TempDir(), "project")
	Generate(root)

	expectedDirs := []string{
		"apps", "services", "libs", "docs", "infrastructure",
		"modules/planner", "modules/generators", "modules/templates",
		"generators/prd_generator", "templates", ".architect", ".codex", ".codex/cache",
	}
	for _, d := range expectedDirs {
		path := filepath.Join(root, d)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("expected dir %q missing: %v", d, err)
		}
	}
}

func TestGenerateExpectedFiles(t *testing.T) {
	root := filepath.Join(t.TempDir(), "project")
	Generate(root)

	expectedFiles := []string{
		"apps/README.md",
		"services/README.md",
		"libs/README.md",
		"infrastructure/README.md",
		"templates/README.md",
		".codex/config.json",
		".codex/memory.md",
	}
	for _, f := range expectedFiles {
		path := filepath.Join(root, f)
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("expected file %q missing: %v", f, err)
			continue
		}
		if len(data) == 0 {
			t.Errorf("file %q is empty", f)
		}
	}
}

func TestGenerateIdempotent(t *testing.T) {
	root := filepath.Join(t.TempDir(), "project")

	r1, err := Generate(root)
	if err != nil {
		t.Fatalf("first Generate: %v", err)
	}

	r2, err := Generate(root)
	if err != nil {
		t.Fatalf("second Generate: %v", err)
	}

	if len(r1.CreatedDirs) != len(r2.CreatedDirs) {
		t.Errorf("dirs count changed: %d vs %d", len(r1.CreatedDirs), len(r2.CreatedDirs))
	}
}
