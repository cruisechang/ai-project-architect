package docphase

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsPhaseBased(t *testing.T) {
	if !IsPhaseBased([]byte("# PRD\n\n## Phase 0\n")) {
		t.Fatal("expected phase heading to be detected")
	}
	if IsPhaseBased([]byte("# PRD\n\n## Goals\n")) {
		t.Fatal("did not expect non-phase doc to be detected as phase-based")
	}
}

func TestCheck(t *testing.T) {
	root := t.TempDir()
	phasedPath := filepath.Join(root, "docs", "PRD.md")
	if err := os.MkdirAll(filepath.Dir(phasedPath), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(phasedPath, []byte("# PRD\n\n## Phase 0\n"), 0o644); err != nil {
		t.Fatalf("write phased file: %v", err)
	}

	status, err := Check(root, "docs/PRD.md")
	if err != nil {
		t.Fatalf("Check phased file: %v", err)
	}
	if !status.Exists || !status.PhaseBased {
		t.Fatalf("unexpected phased status: %+v", status)
	}

	missing, err := Check(root, "docs/SPEC.md")
	if err != nil {
		t.Fatalf("Check missing file: %v", err)
	}
	if missing.Exists || missing.PhaseBased {
		t.Fatalf("unexpected missing status: %+v", missing)
	}
}
