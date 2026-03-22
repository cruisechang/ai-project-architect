package runtime

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"project-generator/internal/architect/model"
)

func TestResolveRoot(t *testing.T) {
	tests := []struct {
		name        string
		parentPath  string
		projectName string
		wantSuffix  string
		wantErr     bool
	}{
		{"current dir with project", ".", "my-project", "my-project", false},
		{"absolute with project", "/tmp", "app", filepath.Join("/tmp", "app"), false},
		{"no project name", "/tmp", "", "/tmp", false},
		{"empty parent defaults", "", "app", "app", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveRoot(tt.parentPath, tt.projectName)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !filepath.IsAbs(got) {
				t.Errorf("result %q is not absolute", got)
			}
			if tt.wantSuffix != "" && !strings.HasSuffix(got, tt.wantSuffix) {
				t.Errorf("result %q does not end with %q", got, tt.wantSuffix)
			}
		})
	}
}

func TestSaveAndLoadContext(t *testing.T) {
	dir := t.TempDir()

	ctx := model.ProjectContext{
		ProjectName:       "round-trip",
		ProjectIdea:       "test save/load",
		ProjectType:       "backend",
		BackendLanguage:   "python",
		BackendFramework:  "fastapi",
		Database:          "mysql",
		Authentication:    "oauth",
		Deployment:        "kubernetes",
		FrontendFramework: "none",
	}

	if err := SaveContext(dir, ctx); err != nil {
		t.Fatalf("SaveContext: %v", err)
	}

	contextFile := filepath.Join(dir, ContextRelPath)
	if _, err := os.Stat(contextFile); err != nil {
		t.Fatalf("context file not created: %v", err)
	}

	loaded, err := LoadContext(dir)
	if err != nil {
		t.Fatalf("LoadContext: %v", err)
	}

	if loaded.ProjectName != ctx.ProjectName {
		t.Errorf("ProjectName = %q, want %q", loaded.ProjectName, ctx.ProjectName)
	}
	if loaded.BackendLanguage != ctx.BackendLanguage {
		t.Errorf("BackendLanguage = %q, want %q", loaded.BackendLanguage, ctx.BackendLanguage)
	}
	if loaded.Database != ctx.Database {
		t.Errorf("Database = %q, want %q", loaded.Database, ctx.Database)
	}
	if loaded.GeneratedAt == "" {
		t.Error("GeneratedAt should be set by EnsureDefaults")
	}
}

func TestLoadContextMissing(t *testing.T) {
	_, err := LoadContext(t.TempDir())
	if err == nil {
		t.Fatal("expected error for missing context file")
	}
}

func TestSaveContextCreatesDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "deep")
	ctx := model.ProjectContext{ProjectName: "test"}

	if err := SaveContext(dir, ctx); err != nil {
		t.Fatalf("SaveContext: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, ContextRelPath)); err != nil {
		t.Fatalf("context file not created: %v", err)
	}
}

func TestBackupIfNeeded(t *testing.T) {
	t.Run("non-existent path", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "nope")
		backup, err := BackupIfNeeded(path, false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if backup != "" {
			t.Errorf("backup = %q, want empty", backup)
		}
	})

	t.Run("existing dir without overwrite", func(t *testing.T) {
		dir := t.TempDir()
		_, err := BackupIfNeeded(dir, false)
		if err == nil {
			t.Fatal("expected error without overwrite")
		}
		if !strings.Contains(err.Error(), "--overwrite") {
			t.Errorf("error = %q, want mention of --overwrite", err.Error())
		}
	})

	t.Run("existing dir with overwrite", func(t *testing.T) {
		parent := t.TempDir()
		target := filepath.Join(parent, "project")
		os.MkdirAll(target, 0o755)
		os.WriteFile(filepath.Join(target, "file.txt"), []byte("data"), 0o644)

		backup, err := BackupIfNeeded(target, true)
		if err != nil {
			t.Fatalf("BackupIfNeeded: %v", err)
		}
		if backup == "" {
			t.Fatal("backup path should not be empty")
		}
		if _, err := os.Stat(target); !os.IsNotExist(err) {
			t.Error("original dir should be renamed")
		}
		if _, err := os.Stat(backup); err != nil {
			t.Errorf("backup dir should exist: %v", err)
		}
	})

	t.Run("existing file not dir", func(t *testing.T) {
		file := filepath.Join(t.TempDir(), "file")
		os.WriteFile(file, []byte("x"), 0o644)
		_, err := BackupIfNeeded(file, true)
		if err == nil {
			t.Fatal("expected error for non-directory")
		}
	})
}
