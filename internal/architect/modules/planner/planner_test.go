package planner

import (
	"testing"
)

func TestInferProjectType(t *testing.T) {
	tests := []struct {
		name     string
		idea     string
		wantType string
	}{
		{"frontend only", "build a react dashboard UI", "frontend"},
		{"backend only", "create a go api service", "backend"},
		{"fullstack explicit", "build a web app with react frontend and go backend", "fullstack"},
		{"no keywords defaults to fullstack", "make a new tool for devs", "fullstack"},
		{"vue frontend", "vue ui for admin", "frontend"},
		{"next frontend", "next app with dashboard", "frontend"},
	}
	p := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := p.Infer(tt.idea, "")
			if ctx.ProjectType != tt.wantType {
				t.Errorf("Infer(%q).ProjectType = %q, want %q", tt.idea, ctx.ProjectType, tt.wantType)
			}
		})
	}
}

func TestInferFrontendFramework(t *testing.T) {
	tests := []struct {
		idea string
		want string
	}{
		{"build a nuxt app", "nuxt"},
		{"vue dashboard", "nuxt"},
		{"next.js app", "next"},
		{"react frontend", "react"},
		{"go api service", "none"},
		{"some project", "react"},
	}
	p := New()
	for _, tt := range tests {
		t.Run(tt.idea, func(t *testing.T) {
			ctx := p.Infer(tt.idea, "")
			if ctx.FrontendFramework != tt.want {
				t.Errorf("FrontendFramework = %q, want %q", ctx.FrontendFramework, tt.want)
			}
		})
	}
}

func TestInferBackendLanguage(t *testing.T) {
	tests := []struct {
		idea         string
		wantLang     string
		wantFramwork string
	}{
		{"golang api", "go", "gin"},
		{"go backend service", "go", "gin"},
		{"python fastapi service", "python", "fastapi"},
		{"python django app", "python", "django"},
		{"node express api", "node", "express"},
		{"nestjs backend", "node", "nestjs"},
		{"react dashboard UI", "none", "none"},
		{"some project", "go", "gin"},
	}
	p := New()
	for _, tt := range tests {
		t.Run(tt.idea, func(t *testing.T) {
			ctx := p.Infer(tt.idea, "")
			if ctx.BackendLanguage != tt.wantLang {
				t.Errorf("BackendLanguage = %q, want %q", ctx.BackendLanguage, tt.wantLang)
			}
			if ctx.BackendFramework != tt.wantFramwork {
				t.Errorf("BackendFramework = %q, want %q", ctx.BackendFramework, tt.wantFramwork)
			}
		})
	}
}

func TestInferDatabase(t *testing.T) {
	tests := []struct {
		idea string
		want string
	}{
		{"api with postgres", "postgresql"},
		{"mysql backend", "mysql"},
		{"sqlite local db", "sqlite"},
		{"mongodb document store", "mongodb"},
		{"react ui only", "none"},
		{"some backend project", "postgresql"},
	}
	p := New()
	for _, tt := range tests {
		t.Run(tt.idea, func(t *testing.T) {
			ctx := p.Infer(tt.idea, "")
			if ctx.Database != tt.want {
				t.Errorf("Database = %q, want %q", ctx.Database, tt.want)
			}
		})
	}
}

func TestInferAuthentication(t *testing.T) {
	tests := []struct {
		idea string
		want string
	}{
		{"api with oauth login", "oauth"},
		{"sso enterprise app", "sso"},
		{"public no auth api", "none"},
		{"some api service", "jwt"},
	}
	p := New()
	for _, tt := range tests {
		t.Run(tt.idea, func(t *testing.T) {
			ctx := p.Infer(tt.idea, "")
			if ctx.Authentication != tt.want {
				t.Errorf("Authentication = %q, want %q", ctx.Authentication, tt.want)
			}
		})
	}
}

func TestInferDeployment(t *testing.T) {
	tests := []struct {
		idea string
		want string
	}{
		{"deploy on kubernetes", "kubernetes"},
		{"k8s cluster", "kubernetes"},
		{"serverless lambda function", "serverless"},
		{"vm based deployment", "vm"},
		{"some project", "docker-compose"},
	}
	p := New()
	for _, tt := range tests {
		t.Run(tt.idea, func(t *testing.T) {
			ctx := p.Infer(tt.idea, "")
			if ctx.Deployment != tt.want {
				t.Errorf("Deployment = %q, want %q", ctx.Deployment, tt.want)
			}
		})
	}
}

func TestNormalizeProjectName(t *testing.T) {
	tests := []struct {
		name     string
		explicit string
		idea     string
		want     string
	}{
		{"explicit name", "My Project", "", "my-project"},
		{"explicit with spaces", "  Hello World  ", "", "hello-world"},
		{"from idea", "", "Build a Go API", "build-a-go-api"},
		{"empty both", "", "", "ai-project-architect-project"},
		{"whitespace only", "   ", "", "ai-project-architect-project"},
		{"special chars", "", "my@cool#project!", "my-cool-project"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeProjectName(tt.explicit, tt.idea)
			if got != tt.want {
				t.Errorf("normalizeProjectName(%q, %q) = %q, want %q", tt.explicit, tt.idea, got, tt.want)
			}
		})
	}
}

func TestToSlug(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Hello World", "hello-world"},
		{"  spaces  ", "spaces"},
		{"special@#$chars!", "special-chars"},
		{"", "ai-project-architect-project"},
		{"---", "ai-project-architect-project"},
		{"a" + string(make([]byte, 100)), ""}, // long string gets truncated
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := toSlug(tt.input)
			if tt.want != "" && got != tt.want {
				t.Errorf("toSlug(%q) = %q, want %q", tt.input, got, tt.want)
			}
			if len(got) > 48 {
				t.Errorf("toSlug result too long: %d chars", len(got))
			}
		})
	}
}
