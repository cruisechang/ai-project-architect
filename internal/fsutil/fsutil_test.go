package fsutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteFile(t *testing.T) {
	t.Run("creates file with parent dirs", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "a", "b", "file.txt")
		data := []byte("hello")

		if err := WriteFile(path, data, 0o644, false); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}

		got, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("ReadFile: %v", err)
		}
		if string(got) != "hello" {
			t.Errorf("content = %q, want %q", got, "hello")
		}
	})

	t.Run("dry run skips write", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "dry.txt")

		if err := WriteFile(path, []byte("data"), 0o644, true); err != nil {
			t.Fatalf("WriteFile dry run: %v", err)
		}

		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Error("file should not exist in dry run")
		}
	})

	t.Run("default permission", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "default.txt")

		if err := WriteFile(path, []byte("x"), 0, false); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}

		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("Stat: %v", err)
		}
		if info.Mode().Perm() != 0o644 {
			t.Errorf("perm = %o, want %o", info.Mode().Perm(), 0o644)
		}
	})
}

func TestEnsureDir(t *testing.T) {
	t.Run("creates nested dirs", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "a", "b", "c")

		if err := EnsureDir(path, false); err != nil {
			t.Fatalf("EnsureDir: %v", err)
		}

		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("dir not created: %v", err)
		}
		if !info.IsDir() {
			t.Error("path is not a directory")
		}
	})

	t.Run("dry run skips creation", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "should-not-exist")

		if err := EnsureDir(path, true); err != nil {
			t.Fatalf("EnsureDir dry run: %v", err)
		}

		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Error("dir should not exist in dry run")
		}
	})
}

func TestExists(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "exists.txt")
	os.WriteFile(file, []byte("x"), 0o644)

	tests := []struct {
		name string
		path string
		want bool
	}{
		{"existing file", file, true},
		{"existing dir", dir, true},
		{"non-existent", filepath.Join(dir, "nope"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Exists(tt.path)
			if err != nil {
				t.Fatalf("Exists: %v", err)
			}
			if got != tt.want {
				t.Errorf("Exists(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestCopyDir(t *testing.T) {
	src := t.TempDir()
	os.MkdirAll(filepath.Join(src, "sub"), 0o755)
	os.WriteFile(filepath.Join(src, "root.txt"), []byte("root"), 0o644)
	os.WriteFile(filepath.Join(src, "sub", "nested.txt"), []byte("nested"), 0o644)

	dst := filepath.Join(t.TempDir(), "copy")

	if err := CopyDir(src, dst); err != nil {
		t.Fatalf("CopyDir: %v", err)
	}

	got, err := os.ReadFile(filepath.Join(dst, "root.txt"))
	if err != nil {
		t.Fatalf("root.txt: %v", err)
	}
	if string(got) != "root" {
		t.Errorf("root.txt = %q", got)
	}

	got, err = os.ReadFile(filepath.Join(dst, "sub", "nested.txt"))
	if err != nil {
		t.Fatalf("nested.txt: %v", err)
	}
	if string(got) != "nested" {
		t.Errorf("nested.txt = %q", got)
	}
}

func TestCopyDirNotADir(t *testing.T) {
	file := filepath.Join(t.TempDir(), "file.txt")
	os.WriteFile(file, []byte("x"), 0o644)

	err := CopyDir(file, t.TempDir())
	if err == nil {
		t.Fatal("expected error for non-directory source")
	}
}
