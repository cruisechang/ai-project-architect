package planner

import (
	"regexp"
	"strings"

	"project-generator/internal/architect/model"
)

type rule struct {
	keywords []string
	value    string
}

var frontendRules = []rule{
	{keywords: []string{"nuxt", "vue"}, value: "nuxt"},
	{keywords: []string{"next"}, value: "next"},
	{keywords: []string{"react"}, value: "react"},
}

var backendRules = []rule{
	{keywords: []string{"golang", "go backend", "go api", "go service", " go "}, value: "go"},
	{keywords: []string{"python", "fastapi", "django"}, value: "python"},
	{keywords: []string{"node", "typescript backend", "javascript backend", "express", "nestjs"}, value: "node"},
}

var backendFrameworkRules = []rule{
	{keywords: []string{"django"}, value: "django"},
	{keywords: []string{"fastapi", "python"}, value: "fastapi"},
	{keywords: []string{"nestjs"}, value: "nestjs"},
	{keywords: []string{"node", "express"}, value: "express"},
	{keywords: []string{"golang", "go backend", "go api", "go service", " go "}, value: "gin"},
}

var databaseRules = []rule{
	{keywords: []string{"postgres", "postgresql"}, value: "postgresql"},
	{keywords: []string{"mysql"}, value: "mysql"},
	{keywords: []string{"sqlite"}, value: "sqlite"},
	{keywords: []string{"mongodb", "mongo"}, value: "mongodb"},
}

var authRules = []rule{
	{keywords: []string{"oauth"}, value: "oauth"},
	{keywords: []string{"sso"}, value: "sso"},
	{keywords: []string{"no auth", "public"}, value: "none"},
}

var deploymentRules = []rule{
	{keywords: []string{"k8s", "kubernetes"}, value: "kubernetes"},
	{keywords: []string{"serverless", "lambda", "cloud function"}, value: "serverless"},
	{keywords: []string{"vm", "server", "bare metal"}, value: "vm"},
}

type Planner struct{}

func New() Planner {
	return Planner{}
}

func (Planner) Infer(idea, explicitProjectName string) model.ProjectContext {
	text := strings.ToLower(strings.TrimSpace(idea))
	ctx := model.ProjectContext{
		ProjectName: normalizeProjectName(explicitProjectName, idea),
		ProjectIdea: strings.TrimSpace(idea),
	}

	hasFrontend := containsAny(text, "frontend", "ui", "react", "vue", "next", "nuxt", "web app", "dashboard")
	hasBackend := containsAny(text, "backend", "api", "service", "microservice", "golang", "go ", "python", "node", "database")

	switch {
	case hasFrontend && hasBackend:
		ctx.ProjectType = "fullstack"
	case hasFrontend:
		ctx.ProjectType = "frontend"
	case hasBackend:
		ctx.ProjectType = "backend"
	default:
		ctx.ProjectType = "fullstack"
	}

	if ctx.ProjectType == "backend" {
		ctx.FrontendFramework = "none"
	} else {
		ctx.FrontendFramework = matchRule(text, frontendRules, "react")
	}

	if ctx.ProjectType == "frontend" {
		ctx.BackendLanguage = "none"
		ctx.BackendFramework = "none"
	} else {
		ctx.BackendLanguage = matchRule(text, backendRules, "go")
		ctx.BackendFramework = matchRule(text, backendFrameworkRules, "gin")
	}

	if ctx.ProjectType == "frontend" {
		ctx.Database = "none"
	} else {
		ctx.Database = matchRule(text, databaseRules, "postgresql")
	}

	ctx.Authentication = matchRule(text, authRules, "jwt")
	ctx.Deployment = matchRule(text, deploymentRules, "docker-compose")

	ctx.EnsureDefaults()
	return ctx
}

// matchRule returns the value of the first rule whose keywords match text,
// or fallback if no rule matches.
func matchRule(text string, rules []rule, fallback string) string {
	for _, r := range rules {
		if containsAny(text, r.keywords...) {
			return r.value
		}
	}
	return fallback
}

func containsAny(text string, keys ...string) bool {
	for _, k := range keys {
		if strings.Contains(text, strings.ToLower(k)) {
			return true
		}
	}
	return false
}

func normalizeProjectName(explicitProjectName, idea string) string {
	if trimmed := strings.TrimSpace(explicitProjectName); trimmed != "" {
		return toSlug(trimmed)
	}
	if strings.TrimSpace(idea) == "" {
		return "ai-project-architect-project"
	}
	return toSlug(idea)
}

func toSlug(value string) string {
	v := strings.ToLower(strings.TrimSpace(value))
	re := regexp.MustCompile(`[^a-z0-9]+`)
	v = re.ReplaceAllString(v, "-")
	v = strings.Trim(v, "-")
	if v == "" {
		return "ai-project-architect-project"
	}
	if len(v) > 48 {
		v = strings.Trim(v[:48], "-")
	}
	if v == "" {
		return "ai-project-architect-project"
	}
	return v
}
