package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExpandPath(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("cannot get home dir: %v", err)
	}

	tests := []struct {
		name    string
		raw     string
		want    string
		wantErr string
	}{
		{"empty", "", "", ""},
		{"whitespace only", "   ", "", ""},
		{"tilde only", "~", home, ""},
		{"tilde with subpath", "~/foo/bar", filepath.Join(home, "foo", "bar"), ""},
		{"relative path", "foo/bar", "", ""}, // will be resolved to abs
		{"unsupported tilde format", "~other/foo", "", "unsupported home path"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExpandPath(tt.raw)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.wantErr)
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("error %q does not contain %q", err.Error(), tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.want != "" && got != tt.want {
				t.Errorf("ExpandPath(%q) = %q, want %q", tt.raw, got, tt.want)
			}
			if tt.want == "" && tt.raw != "" && strings.TrimSpace(tt.raw) != "" {
				if !filepath.IsAbs(got) {
					t.Errorf("ExpandPath(%q) = %q, expected absolute path", tt.raw, got)
				}
			}
		})
	}
}

func TestExpandPathEnvVar(t *testing.T) {
	t.Setenv("TEST_APA_DIR", "/opt/custom")
	got, err := ExpandPath("$TEST_APA_DIR/sub")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "/opt/custom/sub" {
		t.Errorf("ExpandPath with env = %q, want %q", got, "/opt/custom/sub")
	}
}
