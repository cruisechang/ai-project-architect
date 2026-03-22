package generators

import (
	"strings"
	"testing"

	"project-generator/internal/architect/model"
)

func TestDefaultArtifacts(t *testing.T) {
	artifacts := DefaultArtifacts()

	expectedKeys := []string{"prd", "spec", "architecture", "api", "db", "implementation-plan"}
	if len(artifacts) != len(expectedKeys) {
		t.Fatalf("DefaultArtifacts() len = %d, want %d", len(artifacts), len(expectedKeys))
	}
	for i, a := range artifacts {
		if a.Key != expectedKeys[i] {
			t.Errorf("artifact[%d].Key = %q, want %q", i, a.Key, expectedKeys[i])
		}
		if a.FilePath == "" {
			t.Errorf("artifact[%d].FilePath is empty", i)
		}
	}
}

func TestDefaultArtifactsExcludesRepoStructure(t *testing.T) {
	for _, a := range DefaultArtifacts() {
		if a.Key == "repository-structure" {
			t.Error("DefaultArtifacts() should not include repository-structure")
		}
	}
}

func TestResolveTarget(t *testing.T) {
	tests := []struct {
		target   string
		wantLen  int
		wantKeys []string
		wantErr  string
	}{
		{"docs", 6, nil, ""},
		{"all", 6, nil, ""},
		{"", 6, nil, ""},
		{"prd", 1, []string{"prd"}, ""},
		{"spec", 1, []string{"spec"}, ""},
		{"architecture", 1, []string{"architecture"}, ""},
		{"api", 1, []string{"api"}, ""},
		{"db", 1, []string{"db"}, ""},
		{"implementation-plan", 1, []string{"implementation-plan"}, ""},
		{"repository-structure", 1, []string{"repository-structure"}, ""},
		{"project-structure", 1, []string{"repository-structure"}, ""},
		{"scaffold", 0, nil, ""},
		{"invalid-target", 0, nil, "unsupported artifact target"},
	}
	for _, tt := range tests {
		t.Run(tt.target, func(t *testing.T) {
			artifacts, err := ResolveTarget(tt.target)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q", tt.wantErr)
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("error %q missing %q", err.Error(), tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(artifacts) != tt.wantLen {
				t.Fatalf("len = %d, want %d", len(artifacts), tt.wantLen)
			}
			if tt.wantKeys != nil {
				for i, k := range tt.wantKeys {
					if artifacts[i].Key != k {
						t.Errorf("artifact[%d].Key = %q, want %q", i, artifacts[i].Key, k)
					}
				}
			}
		})
	}
}

func TestGenerate(t *testing.T) {
	ctx := model.ProjectContext{
		ProjectName:       "test-gen",
		ProjectIdea:       "test idea",
		ProjectType:       "fullstack",
		FrontendFramework: "react",
		BackendLanguage:   "go",
		BackendFramework:  "gin",
		Database:          "postgresql",
		Authentication:    "jwt",
		Deployment:        "docker-compose",
	}

	for _, a := range orderedArtifacts {
		t.Run(a.Key, func(t *testing.T) {
			result, err := Generate(ctx, a)
			if err != nil {
				t.Fatalf("Generate(%q): %v", a.Key, err)
			}
			if result == "" {
				t.Errorf("Generate(%q) returned empty", a.Key)
			}
		})
	}
}

func TestGenerateUnknownKey(t *testing.T) {
	ctx := model.ProjectContext{}
	_, err := Generate(ctx, Artifact{Key: "unknown"})
	if err == nil {
		t.Fatal("expected error for unknown key")
	}
}
