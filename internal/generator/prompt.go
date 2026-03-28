package generator

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"time"

	templatefs "project-generator/templates"

	"project-generator/internal/config"
)

type TemplateData struct {
	ProjectName        string
	ProjectIdea        string
	ProjectType        string
	AIAgent            string
	IncludeCodex       bool
	IncludeClaude      bool
	ProjectDescription string
	BackendType        string
	FrontendType       string
	DocsType           string
	Architecture       string
	TechStack          string
	UnitTest           string
	APITest            string
	IntegrationTest    string
	E2ETest            string
	PerformanceTest    string
	SecurityTest       string
	UATTest            string
	DockerCompose      string
	PromptTitle        string
	IncludedSkills     []string
	IncludedSkillsCSV  string
	GeneratedAt        string
	Goals              []string
	Requirements       []string
	Notes              []string
}

func buildTemplateData(opts config.CreateOptions) TemplateData {
	goals, requirements, notes := defaultsByProjectType(opts.ProjectType)
	skillsCSV := "none"
	if len(opts.SelectedSkills) > 0 {
		skillsCSV = strings.Join(opts.SelectedSkills, ", ")
	}
	description := opts.Description
	if description == "" {
		description = summarizeIdea(opts.Idea)
	}
	if description == "" {
		description = "TBD"
	}
	architecture := opts.Architecture
	if architecture == "" {
		architecture = "TBD"
	}
	techStack := opts.TechStack
	if techStack == "" {
		techStack = "TBD"
	}

	return TemplateData{
		ProjectName:        opts.Name,
		ProjectIdea:        opts.Idea,
		ProjectType:        opts.ProjectType,
		AIAgent:            opts.AIAgent,
		IncludeCodex:       opts.AIAgent == "codex" || opts.AIAgent == "universal",
		IncludeClaude:      opts.AIAgent == "claude-code" || opts.AIAgent == "universal",
		ProjectDescription: description,
		BackendType:        opts.BackendType,
		FrontendType:       opts.FrontendType,
		DocsType:           opts.DocsType,
		Architecture:       architecture,
		TechStack:          techStack,
		UnitTest:           yesNoLabel(opts.UnitTest),
		APITest:            yesNoLabel(opts.APITest),
		IntegrationTest:    yesNoLabel(opts.IntegrationTest),
		E2ETest:            yesNoLabel(opts.E2ETest),
		PerformanceTest:    yesNoLabel(opts.PerformanceTest),
		SecurityTest:       yesNoLabel(opts.SecurityTest),
		UATTest:            yesNoLabel(opts.UATTest),
		DockerCompose:      yesNoLabel(opts.DockerCompose),
		PromptTitle:        opts.PromptTitle,
		IncludedSkills:     opts.SelectedSkills,
		IncludedSkillsCSV:  skillsCSV,
		GeneratedAt:        time.Now().Format(time.RFC3339),
		Goals:              goals,
		Requirements:       requirements,
		Notes:              notes,
	}
}

func summarizeIdea(idea string) string {
	trimmed := strings.Join(strings.Fields(strings.TrimSpace(idea)), " ")
	if trimmed == "" {
		return ""
	}
	runes := []rune(trimmed)
	const maxLen = 140
	if len(runes) <= maxLen {
		return trimmed
	}
	return strings.TrimSpace(string(runes[:maxLen-1])) + "..."
}

func renderTemplate(templateName string, data TemplateData) (string, error) {
	raw, err := templatefs.FS.ReadFile(templateName)
	if err != nil {
		return "", fmt.Errorf("cannot load template %s: %w", templateName, err)
	}

	tmpl, err := template.New(templateName).Parse(string(raw))
	if err != nil {
		return "", fmt.Errorf("cannot parse template %s: %w", templateName, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("cannot execute template %s: %w", templateName, err)
	}
	return buf.String(), nil
}

func defaultsByProjectType(projectType string) ([]string, []string, []string) {
	switch projectType {
	case "frontend-only":
		return []string{
				"Deliver a polished frontend product with coherent UX and maintainable client architecture.",
				"Keep critical user journeys explicit, testable, and observable from day one.",
			}, []string{
				"Define UI scope and key screens for the first release.",
				"Define component boundaries and state/data flow strategy.",
				"Define routing, error handling, and loading states.",
				"Define frontend test strategy (unit/component/e2e).",
				"Define asset/build pipeline and deployment considerations.",
			}, []string{
				"Prioritize responsive behavior and production performance budgets.",
			}
	case "backend-only":
		return []string{
				"Ship a reliable backend with clear contracts, operational safety, and maintainable architecture.",
				"Keep service boundaries and failure handling explicit and testable.",
			}, []string{
				"Define API/service contracts and versioning policy.",
				"Define data model and persistence strategy.",
				"Define authentication and authorization approach.",
				"Define background jobs/async processing where needed.",
				"Define observability (logs/metrics/traces) and SLOs.",
				"Define rollback and incident response workflow.",
			}, []string{
				"Prefer stable interfaces and explicit error semantics.",
			}
	case "full-stack", "ai-app", "web-app", "devops-tool", "internal-tool", "platform-service":
		return []string{
				"Deliver a complete product experience across frontend and backend.",
				"Keep integration boundaries clear while preserving fast delivery cycles.",
			}, []string{
				"Define UI scope, user journeys, and API contracts together.",
				"Define authentication/session model end-to-end.",
				"Define backend service boundaries and data ownership.",
				"Define integration and e2e test coverage for critical flows.",
				"Define deployment topology and release/rollback strategy.",
			}, []string{
				"Validate behavior across UI/API/Data boundaries before each release.",
			}
	default:
		return []string{"Define clear project outcomes."}, []string{"Document technical decisions."}, []string{"Keep scope focused."}
	}
}

func yesNoLabel(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "yes":
		return "Yes"
	case "no":
		return "No"
	default:
		return "TBD"
	}
}
