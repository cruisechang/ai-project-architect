package apa

import (
	"testing"

	"project-generator/internal/architect/model"
	"project-generator/internal/config"
)

func TestMapContextToOpts_FillsMissingFields(t *testing.T) {
	ctx := model.ProjectContext{
		ProjectType:       "fullstack",
		BackendLanguage:   "go",
		FrontendFramework: "react",
		Deployment:        "docker-compose",
	}
	opts := config.CreateOptions{}
	mapContextToOpts(ctx, &opts)

	if opts.BackendType != "go" {
		t.Errorf("BackendType: want go, got %q", opts.BackendType)
	}
	if opts.FrontendType != "react" {
		t.Errorf("FrontendType: want react, got %q", opts.FrontendType)
	}
	if opts.Architecture != "frontend-backend" {
		t.Errorf("Architecture: want frontend-backend, got %q", opts.Architecture)
	}
	if opts.DockerCompose != "yes" {
		t.Errorf("DockerCompose: want yes, got %q", opts.DockerCompose)
	}
	if opts.ProjectType != "internal-tool" {
		t.Errorf("ProjectType: want internal-tool, got %q", opts.ProjectType)
	}
	if opts.TechStack == "" {
		t.Error("TechStack: want non-empty")
	}
}

func TestMapContextToOpts_ExplicitFlagsTakePriority(t *testing.T) {
	ctx := model.ProjectContext{
		ProjectType:       "fullstack",
		BackendLanguage:   "go",
		FrontendFramework: "react",
	}
	opts := config.CreateOptions{
		BackendType:  "python", // explicit flag should not be overridden
		Architecture: "cli-tool",
	}
	mapContextToOpts(ctx, &opts)

	if opts.BackendType != "python" {
		t.Errorf("BackendType: explicit flag should not be overridden, got %q", opts.BackendType)
	}
	if opts.Architecture != "cli-tool" {
		t.Errorf("Architecture: explicit flag should not be overridden, got %q", opts.Architecture)
	}
	// FrontendType was empty, so inferred value is used
	if opts.FrontendType != "react" {
		t.Errorf("FrontendType: want react (inferred), got %q", opts.FrontendType)
	}
}

func TestInferArchitectureFromContext(t *testing.T) {
	cases := []struct {
		projectType string
		frontend    string
		want        string
	}{
		{"backend", "", "backend-service"},
		{"frontend", "react", "frontend-app"},
		{"fullstack", "next", "fullstack-web-app"},
		{"fullstack", "nuxt", "fullstack-web-app"},
		{"fullstack", "react", "frontend-backend"},
		{"fullstack", "vue", "frontend-backend"},
		{"", "", ""},
	}
	for _, c := range cases {
		ctx := model.ProjectContext{
			ProjectType:       c.projectType,
			FrontendFramework: c.frontend,
		}
		got := inferArchitectureFromContext(ctx)
		if got != c.want {
			t.Errorf("inferArchitectureFromContext(%q, frontend=%q): want %q, got %q",
				c.projectType, c.frontend, c.want, got)
		}
	}
}

func TestBuildTechStack(t *testing.T) {
	cases := []struct {
		arch     string
		backend  string
		frontend string
		want     string
	}{
		{"frontend-backend", "go", "react", "frontend-backend | go | react"},
		{"backend-service", "go", "none", "backend-service | go"},
		{"frontend-app", "none", "react", "frontend-app | react"},
		{"", "", "", ""},
		{"fullstack-web-app", "none", "next", "fullstack-web-app | next"},
	}
	for _, c := range cases {
		got := buildTechStack(c.arch, c.backend, c.frontend)
		if got != c.want {
			t.Errorf("buildTechStack(%q, %q, %q): want %q, got %q",
				c.arch, c.backend, c.frontend, c.want, got)
		}
	}
}

func TestMapContextToOpts_DockerComposeInference(t *testing.T) {
	cases := []struct {
		deployment string
		want       string
	}{
		{"docker-compose", "yes"},
		{"docker compose", "yes"},
		{"kubernetes", "no"},
		{"serverless", "no"},
		{"vm", "no"},
	}
	for _, c := range cases {
		ctx := model.ProjectContext{Deployment: c.deployment}
		opts := config.CreateOptions{}
		mapContextToOpts(ctx, &opts)
		if opts.DockerCompose != c.want {
			t.Errorf("deployment=%q: DockerCompose want %q, got %q",
				c.deployment, c.want, opts.DockerCompose)
		}
	}
}
