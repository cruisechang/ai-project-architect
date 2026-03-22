package generators_test

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"project-generator/internal/architect/model"
	generators "project-generator/internal/architect/modules/generators"
)

// Run with -update to regenerate golden files.
var update = flag.Bool("update", false, "update golden files")

// goldenCtx is a fixed context used for all golden tests.
var goldenCtx = model.ProjectContext{
	ProjectName:       "golden-project",
	ProjectIdea:       "a fullstack web app with react and go backend",
	ProjectType:       "fullstack",
	FrontendFramework: "react",
	BackendLanguage:   "go",
	BackendFramework:  "gin",
	Database:          "postgresql",
	Authentication:    "jwt",
	Deployment:        "docker-compose",
}

func TestGeneratorsGolden(t *testing.T) {
	targets := []string{"prd", "spec", "architecture", "api", "db", "implementation-plan", "repository-structure"}

	for _, target := range targets {
		t.Run(target, func(t *testing.T) {
			artifacts, err := generators.ResolveTarget(target)
			if err != nil {
				t.Fatalf("ResolveTarget(%q): %v", target, err)
			}
			if len(artifacts) != 1 {
				t.Fatalf("expected 1 artifact for %q, got %d", target, len(artifacts))
			}

			got, err := generators.Generate(goldenCtx, artifacts[0])
			if err != nil {
				t.Fatalf("Generate(%q): %v", target, err)
			}
			if got == "" {
				t.Fatalf("Generate(%q) returned empty", target)
			}

			goldenPath := filepath.Join("testdata", target+".golden.md")
			if *update {
				if err := os.WriteFile(goldenPath, []byte(got), 0o644); err != nil {
					t.Fatalf("write golden %s: %v", goldenPath, err)
				}
				t.Logf("updated golden: %s", goldenPath)
				return
			}

			want, err := os.ReadFile(goldenPath)
			if err != nil {
				t.Fatalf("golden file missing: %s (run with -update to create)", goldenPath)
			}
			if got != string(want) {
				t.Errorf("output mismatch for %q\n--- want ---\n%s\n--- got ---\n%s", target, want, got)
			}
		})
	}
}
