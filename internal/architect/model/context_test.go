package model

import (
	"testing"
)

func TestEnsureDefaults(t *testing.T) {
	tests := []struct {
		name  string
		input ProjectContext
		check func(t *testing.T, c ProjectContext)
	}{
		{
			name:  "empty context gets fullstack defaults",
			input: ProjectContext{},
			check: func(t *testing.T, c ProjectContext) {
				if c.ProjectType != "fullstack" {
					t.Errorf("ProjectType = %q", c.ProjectType)
				}
				if c.FrontendFramework != "react" {
					t.Errorf("FrontendFramework = %q", c.FrontendFramework)
				}
				if c.BackendLanguage != "go" {
					t.Errorf("BackendLanguage = %q", c.BackendLanguage)
				}
				if c.BackendFramework != "gin" {
					t.Errorf("BackendFramework = %q", c.BackendFramework)
				}
				if c.Database != "postgresql" {
					t.Errorf("Database = %q", c.Database)
				}
				if c.Authentication != "jwt" {
					t.Errorf("Authentication = %q", c.Authentication)
				}
				if c.Deployment != "docker-compose" {
					t.Errorf("Deployment = %q", c.Deployment)
				}
				if c.GeneratedAt == "" {
					t.Error("GeneratedAt should be set")
				}
			},
		},
		{
			name:  "backend project gets no frontend",
			input: ProjectContext{ProjectType: "backend"},
			check: func(t *testing.T, c ProjectContext) {
				if c.FrontendFramework != "none" {
					t.Errorf("FrontendFramework = %q, want none", c.FrontendFramework)
				}
			},
		},
		{
			name:  "frontend project gets no backend",
			input: ProjectContext{ProjectType: "frontend"},
			check: func(t *testing.T, c ProjectContext) {
				if c.BackendLanguage != "none" {
					t.Errorf("BackendLanguage = %q, want none", c.BackendLanguage)
				}
				if c.Database != "none" {
					t.Errorf("Database = %q, want none", c.Database)
				}
			},
		},
		{
			name:  "python gets fastapi",
			input: ProjectContext{BackendLanguage: "python"},
			check: func(t *testing.T, c ProjectContext) {
				if c.BackendFramework != "fastapi" {
					t.Errorf("BackendFramework = %q, want fastapi", c.BackendFramework)
				}
			},
		},
		{
			name:  "node gets express",
			input: ProjectContext{BackendLanguage: "node"},
			check: func(t *testing.T, c ProjectContext) {
				if c.BackendFramework != "express" {
					t.Errorf("BackendFramework = %q, want express", c.BackendFramework)
				}
			},
		},
		{
			name:  "none backend gets none framework",
			input: ProjectContext{BackendLanguage: "none"},
			check: func(t *testing.T, c ProjectContext) {
				if c.BackendFramework != "none" {
					t.Errorf("BackendFramework = %q, want none", c.BackendFramework)
				}
			},
		},
		{
			name:  "does not overwrite existing values",
			input: ProjectContext{ProjectType: "backend", Database: "mysql", Authentication: "oauth", Deployment: "kubernetes", GeneratedAt: "fixed"},
			check: func(t *testing.T, c ProjectContext) {
				if c.Database != "mysql" {
					t.Errorf("Database = %q", c.Database)
				}
				if c.Authentication != "oauth" {
					t.Errorf("Authentication = %q", c.Authentication)
				}
				if c.Deployment != "kubernetes" {
					t.Errorf("Deployment = %q", c.Deployment)
				}
				if c.GeneratedAt != "fixed" {
					t.Errorf("GeneratedAt = %q", c.GeneratedAt)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.input
			c.EnsureDefaults()
			tt.check(t, c)
		})
	}
}
