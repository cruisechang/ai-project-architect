package apa

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolvePromptReviewer_DocsOnlySkipsReviewer(t *testing.T) {
	got, err := resolvePromptReviewer(false, "", true, true, bufio.NewReader(strings.NewReader("apa-codex-review\n")))
	if err != nil {
		t.Fatalf("resolvePromptReviewer returned error: %v", err)
	}
	if got != "" {
		t.Fatalf("expected empty reviewer in docs-only mode, got %q", got)
	}
}

func TestResolvePromptReviewer_ExplicitReviewer(t *testing.T) {
	got, err := resolvePromptReviewer(true, "apa-codex-review", false, false, nil)
	if err != nil {
		t.Fatalf("resolvePromptReviewer returned error: %v", err)
	}
	if got != "apa-codex-review" {
		t.Fatalf("expected explicit reviewer, got %q", got)
	}
}

func TestResolvePromptReviewer_InteractivePrompt(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader("apa-claude-review\n"))
	got, err := resolvePromptReviewer(false, "", false, true, reader)
	if err != nil {
		t.Fatalf("resolvePromptReviewer returned error: %v", err)
	}
	if got != "apa-claude-review" {
		t.Fatalf("expected prompted reviewer, got %q", got)
	}
}

func TestResolvePromptReviewer_NonInteractiveDefaultsAgentSelf(t *testing.T) {
	got, err := resolvePromptReviewer(false, "", false, false, nil)
	if err != nil {
		t.Fatalf("resolvePromptReviewer returned error: %v", err)
	}
	if got != "agent-self" {
		t.Fatalf("expected agent-self default, got %q", got)
	}
}

func TestResolvePromptReviewer_InvalidReviewer(t *testing.T) {
	_, err := resolvePromptReviewer(true, "bad-reviewer", false, false, nil)
	if err == nil {
		t.Fatal("expected invalid reviewer error")
	}
}

func TestPrintPromptOutput_IncludesReviewer(t *testing.T) {
	root := t.TempDir()
	mustWritePromptFile(t, filepath.Join(root, "docs", "PRD.md"), "# PRD\n\n## Phase 0\n")

	output := captureStdout(t, func() {
		printPromptOutput(root, false, "apa-codex-review")
	})

	if !strings.Contains(output, "Reviewer: apa-codex-review") {
		t.Fatalf("expected reviewer label in prompt output, got:\n%s", output)
	}
	if !strings.Contains(output, "Use reviewer `apa-codex-review` for every round unless I explicitly switch it.") {
		t.Fatalf("expected reviewer workflow instruction in prompt output, got:\n%s", output)
	}
}

func mustWritePromptFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	original := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = w
	done := make(chan string, 1)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		done <- buf.String()
	}()

	fn()

	_ = w.Close()
	os.Stdout = original
	return <-done
}
