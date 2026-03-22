package templates

import (
	"strings"
	"testing"

	"project-generator/internal/architect/model"
)

func TestRenderAllTemplates(t *testing.T) {
	ctx := model.ProjectContext{
		ProjectName:       "test-project",
		ProjectIdea:       "a test idea",
		ProjectType:       "fullstack",
		FrontendFramework: "react",
		BackendLanguage:   "go",
		BackendFramework:  "gin",
		Database:          "postgresql",
		Authentication:    "jwt",
		Deployment:        "docker-compose",
	}

	for name := range catalog {
		t.Run(name, func(t *testing.T) {
			result, err := Render(name, ctx)
			if err != nil {
				t.Fatalf("Render(%q): %v", name, err)
			}
			if result == "" {
				t.Errorf("Render(%q) returned empty", name)
			}
		})
	}
}

func TestRenderSubstitution(t *testing.T) {
	ctx := model.ProjectContext{
		ProjectName:       "my-app",
		ProjectIdea:       "build something cool",
		ProjectType:       "backend",
		BackendLanguage:   "python",
		BackendFramework:  "fastapi",
		Database:          "mysql",
		Authentication:    "oauth",
		Deployment:        "kubernetes",
		FrontendFramework: "none",
	}

	tests := []struct {
		template string
		contains []string
	}{
		{"PRD", []string{"my-app", "build something cool", "backend"}},
		{"SPEC", []string{"my-app", "backend", "python", "fastapi", "mysql"}},
		{"ARCHITECTURE", []string{"kubernetes"}},
		{"API", []string{"oauth"}},
		{"DB_SCHEMA", []string{"mysql"}},
		{"REPOSITORY_STRUCTURE", []string{"my-app"}},
	}
	for _, tt := range tests {
		t.Run(tt.template, func(t *testing.T) {
			result, err := Render(tt.template, ctx)
			if err != nil {
				t.Fatalf("Render(%q): %v", tt.template, err)
			}
			for _, s := range tt.contains {
				if !strings.Contains(result, s) {
					t.Errorf("Render(%q) missing %q", tt.template, s)
				}
			}
		})
	}
}

func TestRenderUnknownTemplate(t *testing.T) {
	ctx := model.ProjectContext{}
	_, err := Render("NONEXISTENT", ctx)
	if err == nil {
		t.Fatal("expected error for unknown template")
	}
	if !strings.Contains(err.Error(), "unknown template") {
		t.Errorf("error = %q, want 'unknown template'", err.Error())
	}
}
