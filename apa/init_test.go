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
	if opts.Architecture != "web-app-server" {
		t.Errorf("Architecture: want web-app-server, got %q", opts.Architecture)
	}
	if opts.DockerCompose != "yes" {
		t.Errorf("DockerCompose: want yes, got %q", opts.DockerCompose)
	}
	if opts.ProjectType != "full-stack" {
		t.Errorf("ProjectType: want full-stack, got %q", opts.ProjectType)
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
		Architecture: "cli",
	}
	mapContextToOpts(ctx, &opts)

	if opts.BackendType != "python" {
		t.Errorf("BackendType: explicit flag should not be overridden, got %q", opts.BackendType)
	}
	if opts.Architecture != "cli" {
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
		{"backend", "", "server"},
		{"frontend", "react", "web-app"},
		{"fullstack", "next", "web-app-server"},
		{"fullstack", "nuxt", "web-app-server"},
		{"fullstack", "react", "web-app-server"},
		{"fullstack", "vue", "web-app-server"},
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
		{"web-app-server", "go", "react", "web-app-server | go | react"},
		{"server", "go", "none", "server | go"},
		{"web-app", "none", "react", "web-app | react"},
		{"", "", "", ""},
		{"web-app-server", "none", "next", "web-app-server | next"},
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
