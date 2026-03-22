package model

import "time"

// ProjectContext is the deterministic planning model used by the architect pipeline.
type ProjectContext struct {
	ProjectName       string `json:"project_name"`
	ProjectIdea       string `json:"project_idea"`
	ProjectType       string `json:"project_type"` // frontend | backend | fullstack
	FrontendFramework string `json:"frontend_framework"`
	BackendLanguage   string `json:"backend_language"`
	BackendFramework  string `json:"backend_framework"`
	Database          string `json:"database"`
	Authentication    string `json:"authentication"`
	Deployment        string `json:"deployment"`
	GeneratedAt       string `json:"generated_at"`
}

func (c *ProjectContext) EnsureDefaults() {
	if c.ProjectType == "" {
		c.ProjectType = "fullstack"
	}
	if c.FrontendFramework == "" {
		if c.ProjectType == "backend" {
			c.FrontendFramework = "none"
		} else {
			c.FrontendFramework = "react"
		}
	}
	if c.BackendLanguage == "" {
		if c.ProjectType == "frontend" {
			c.BackendLanguage = "none"
		} else {
			c.BackendLanguage = "go"
		}
	}
	if c.BackendFramework == "" {
		switch c.BackendLanguage {
		case "go":
			c.BackendFramework = "gin"
		case "python":
			c.BackendFramework = "fastapi"
		case "node":
			c.BackendFramework = "express"
		default:
			c.BackendFramework = "none"
		}
	}
	if c.Database == "" {
		if c.ProjectType == "frontend" {
			c.Database = "none"
		} else {
			c.Database = "postgresql"
		}
	}
	if c.Authentication == "" {
		c.Authentication = "jwt"
	}
	if c.Deployment == "" {
		c.Deployment = "docker-compose"
	}
	if c.GeneratedAt == "" {
		c.GeneratedAt = time.Now().UTC().Format(time.RFC3339)
	}
}
